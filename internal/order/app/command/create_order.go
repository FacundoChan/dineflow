package command

import (
	"context"
	"github.com/FacundoChan/gorder-v1/common/decorator"
	"github.com/FacundoChan/gorder-v1/common/genproto/orderpb"
	domain "github.com/FacundoChan/gorder-v1/order/domain/order"
	"github.com/sirupsen/logrus"
)

type CreateOrder struct {
	CustomerID string
	Items      []*orderpb.ItemWithQuantity
}

type CreateOrderResult struct {
	OrderID string
}

type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]

type createOrderHandler struct {
	orderRepo domain.Repository
	//stockGRPC
}

func NewCreateOrderHandler(orderRepo domain.Repository, logger *logrus.Entry, metricsClient decorator.MetricsClient) CreateOrderHandler {
	if orderRepo == nil {
		panic("orderRepo is nil")
	}
	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
		createOrderHandler{orderRepo: orderRepo},
		logger,
		metricsClient,
	)

}

func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
	// TODO: call stock grpc to get items
	var stockResponse []*orderpb.Item
	for _, item := range cmd.Items {
		stockResponse = append(stockResponse, &orderpb.Item{
			ID: item.ID,
			//Name:     "",
			Quantity: item.Quantity,
			//PriceID:  "",
		})
	}
	order, err := c.orderRepo.Create(ctx, &domain.Order{
		ID:          "",
		CustomerID:  cmd.CustomerID,
		Status:      "",
		PaymentLink: "",
		Items:       stockResponse,
	})

	if err != nil {
		return nil, err
	}
	return &CreateOrderResult{OrderID: order.ID}, nil
}
