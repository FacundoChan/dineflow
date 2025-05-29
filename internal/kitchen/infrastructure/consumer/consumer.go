package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/FacundoChan/dineflow/common/broker"
	"github.com/FacundoChan/dineflow/common/genproto/orderpb"
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
	logrus.Infof("Received message from %s: %s", q.Name, msg.Body)

	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
	tr := otel.Tracer("rabbit-mq")
	mqCtx, span := tr.Start(ctx, fmt.Sprintf("rabbit-mq.%s.consume", q.Name))

	var err error

	defer func() {
		span.End()
		if err != nil {
			_ = msg.Nack(false, false)
		} else {
			_ = msg.Ack(false)
		}
	}()

	order := &Order{}
	if err := json.Unmarshal(msg.Body, order); err != nil {
		logrus.Infof("failed to unmarshal msg to order, err: %v", err)
		return
	}

	if order.Status != "paid" {
		err = errors.New("order not paid, can not be cooked")
		return
	}

	cook(order)

	span.AddEvent(fmt.Sprintf("order_cook: %v", order.ID))
	if err := c.orderGPRC.UpdateOrder(mqCtx, &orderpb.Order{
		ID:          order.ID,
		CustomerID:  order.CustomerID,
		Status:      "ready",
		Items:       order.Items,
		PaymentLink: order.PaymentLink,
	}); err != nil {
		if err = broker.HandleRetry(mqCtx, ch, &msg); err != nil {
			logrus.Warnf("kitchen: error handling retry: err=%v", err)
		}
		return
	}

	span.AddEvent("kitchen.order.finished.updated")
	logrus.Info("successfully consumed")
}

func cook(order *Order) {
	logrus.Printf("cooking order: %s", order.ID)
	time.Sleep(5 * time.Second)
	logrus.Printf("order %s done!", order.ID)
}
