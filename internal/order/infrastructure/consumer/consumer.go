package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/otel"

	"github.com/FacundoChan/gorder-v1/common/broker"
	"github.com/FacundoChan/gorder-v1/order/app"
	"github.com/FacundoChan/gorder-v1/order/app/command"
	domain "github.com/FacundoChan/gorder-v1/order/domain/order"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	app app.Application
}

func NewConsumer(app app.Application) *Consumer {
	return &Consumer{app: app}
}

func (c *Consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(broker.EventOrderPaid, true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	err = ch.QueueBind(q.Name, "", broker.EventOrderPaid, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
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
	logrus.Infof("kitchen receive a msg from %s=%v", q.Name, string(msg.Body))
	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
	t := otel.Tracer("rabbit-mq")
	_, span := t.Start(ctx, fmt.Sprintf("rabbit-mq.%s.consume", q.Name))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			_ = msg.Nack(false, false)
		} else {
			_ = msg.Ack(false)
		}
	}()

	order := &domain.Order{}
	if err := json.Unmarshal(msg.Body, order); err != nil {
		logrus.Infof("unmarshal order error: %s", err)
		return
	}

	_, err = c.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
		Order: order,
		UpdateFunc: func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
			if err := order.IsPaid(); err != nil {
				return nil, err
			}
			return order, nil
		},
	})
	if err != nil {
		logrus.Infof("update order error: %s, orderID=%s", err, order.ID)
		if err = broker.HandleRetry(ctx, ch, &msg); err != nil {
			logrus.Warnf("retry_error handling retry, msgID=%s, err=%v", msg.MessageId, err)
		}
		return
	}

	span.AddEvent("order.updated")
	logrus.Info("order updated")
}
