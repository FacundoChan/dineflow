package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/FacundoChan/dineflow/common/broker"
	"github.com/FacundoChan/dineflow/common/genproto/orderpb"
	"github.com/FacundoChan/dineflow/common/logging"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

type Order struct {
	ID          string
	CustomerID  string
	Status      string
	PaymentLink string
	Items       []*orderpb.Item
}

type OrderService interface {
	UpdateOrder(ctx context.Context, request *orderpb.Order) error
}

type Consumer struct {
	orderGPRC OrderService
}

func NewConsumer(orderGRPC OrderService) *Consumer {
	return &Consumer{
		orderGPRC: orderGRPC,
	}
}

func (c *Consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare("", true, false, true, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	if err = ch.QueueBind(q.Name, "", broker.EventOrderPaid, false, nil); err != nil {
		logrus.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Warnf("Failed to consume message: queue=%s, err=%s", q.Name, err)
	}

	forever := make(chan struct{})
	go func() {
		for msg := range msgs {
			c.handleMessage(msg, q, ch)
		}
	}()

	<-forever
}

func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue, ch *amqp.Channel) {
	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
	t := otel.Tracer("rabbit-mq")
	_, span := t.Start(ctx, fmt.Sprintf("rabbit-mq.%s.consume", q.Name))
	defer span.End()

	logging.Infof(ctx, nil, "kitchen received a message from %s, msg=%v", q.Name, string(msg.Body))

	var err error
	defer func() {
		if err != nil {
			logging.Warnf(ctx, nil, "failed to consumed message from %v, msg=%v err: %v", q.Name, msg, err)
			_ = msg.Nack(false, false)
		} else {
			logging.Infof(ctx, nil, "comsumed message successfully from %v, msg=%v", q.Name, msg)
			_ = msg.Ack(false)
		}
	}()

	order := &Order{}
	if unmarshalErr := json.Unmarshal(msg.Body, order); unmarshalErr != nil {
		err = errors.Wrap(unmarshalErr, "failed to unmarshal msg to order")
		return
	}

	if order.Status != "paid" {
		err = errors.New("order not paid, can not be cooked")
		return
	}

	cook(ctx, order)

	span.AddEvent(fmt.Sprintf("order_cook: %v", order.ID))
	if updateErr := c.orderGPRC.UpdateOrder(ctx, &orderpb.Order{
		ID:          order.ID,
		CustomerID:  order.CustomerID,
		Status:      "ready",
		Items:       order.Items,
		PaymentLink: order.PaymentLink,
	}); updateErr != nil {
		logging.Errorf(ctx, nil, "error updating order: orderID=%s, err=%s", order.ID, updateErr)
		if retryErr := broker.HandleRetry(ctx, ch, &msg); retryErr != nil {
			err = errors.Wrapf(retryErr, "kitchen: error handling retry, msgID=%s, err=%v", msg.MessageId, retryErr)
		} else {
			err = updateErr
		}
		return
	}

	span.AddEvent("kitchen.order.finished.updated")
}

func cook(ctx context.Context, order *Order) {
	logrus.WithContext(ctx).Printf("cooking order: %s", order.ID)
	time.Sleep(5 * time.Second)
	logrus.WithContext(ctx).Printf("order %s done!", order.ID)
}
