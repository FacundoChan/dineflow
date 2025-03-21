package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/FacundoChan/gorder-v1/common/broker"
	"github.com/FacundoChan/gorder-v1/common/decorator"
	"github.com/FacundoChan/gorder-v1/order/app/query"
	"github.com/FacundoChan/gorder-v1/order/convertor"
	domain "github.com/FacundoChan/gorder-v1/order/domain/order"
	"github.com/FacundoChan/gorder-v1/order/entity"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

type CreateOrder struct {
	CustomerID string
	Items      []*entity.ItemWithQuantity
}

type CreateOrderResult struct {
	OrderID string
}

type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]

type createOrderHandler struct {
	orderRepo domain.Repository
	stockGRPC query.StockService
	channel   *amqp.Channel
}

func NewCreateOrderHandler(
	orderRepo domain.Repository,
	stockGRPC query.StockService,
	channel *amqp.Channel,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient) CreateOrderHandler {
	if orderRepo == nil {
		panic("orderRepo is nil")
	}
	if stockGRPC == nil {
		logrus.Panic("stockGRPC is nil")
	}
	if channel == nil {
		logrus.Panic("channel is nil")
	}
	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
		createOrderHandler{
			orderRepo: orderRepo,
			stockGRPC: stockGRPC,
			channel:   channel,
		},
		logger,
		metricsClient,
	)

}

func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
	q, err := c.channel.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	t := otel.Tracer("rabbit-mq")
	ctx, span := t.Start(ctx, fmt.Sprintf("rabbit-mq.%s.publish", q.Name))
	defer span.End()

	validItems, err := c.validate(ctx, cmd.Items)
	if err != nil {
		logrus.WithError(err).Error("invalid items")
	}

	// HACK: should be updated
	var validItemsStrings []string
	for _, item := range validItems {
		validItemsStrings = append(validItemsStrings, item.ID)
	}
	items, err := c.stockGRPC.GetItems(ctx, validItemsStrings)
	if err != nil {
		logrus.Error(err)
	}

	for _, item := range items {
		for _, validItem := range validItems {
			if item.ID == validItem.ID {
				validItem.Name = item.Name
				validItem.PriceID = item.PriceID
				break
			}
		}
	}

	order, err := c.orderRepo.Create(ctx, &domain.Order{
		ID:          "",
		CustomerID:  cmd.CustomerID,
		Status:      "",
		PaymentLink: "",
		Items:       validItems,
	})

	if err != nil {
		return nil, err
	}

	marshalledOrder, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}
	headers := broker.InjectRabbitMQHeaders(ctx)
	err = c.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         marshalledOrder,
		Headers:      headers,
	})
	if err != nil {
		return nil, err
	}

	return &CreateOrderResult{OrderID: order.ID}, nil
}

func (c createOrderHandler) validate(ctx context.Context, items []*entity.ItemWithQuantity) ([]*entity.Item, error) {
	if len(items) == 0 {
		return nil, errors.New("at least one item is required")
	}

	items = packItems(items)
	response, err := c.stockGRPC.CheckIfItemsInStock(ctx, convertor.NewItemWithQuantityConvertor().EntitiesToProtos(items))
	if err != nil {
		return nil, err
	}

	return convertor.NewItemConvertor().ProtosToEntities(response.Items), nil
	//var ids []string
	//for _, i := range items {
	//	ids = append(ids, i.ID)
	//}
	//
	//return c.stockGRPC.GetItems(ctx, ids)
}

func packItems(items []*entity.ItemWithQuantity) []*entity.ItemWithQuantity {
	merged := make(map[string]int32)
	for _, item := range items {
		merged[item.ID] += item.Quantity
	}
	for id, quantity := range merged {
		logrus.Debugf("merged item %v with quantity: %d", id, quantity)

	}
	var res []*entity.ItemWithQuantity
	for id, quantity := range merged {
		res = append(res, &entity.ItemWithQuantity{
			ID:       id,
			Quantity: quantity,
		})
	}
	return res
}
