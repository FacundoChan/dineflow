package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/FacundoChan/dineflow/common/broker"
	"github.com/FacundoChan/dineflow/common/genproto/orderpb"
	"github.com/FacundoChan/dineflow/common/logging"
	"github.com/FacundoChan/dineflow/payment/app"
	"github.com/FacundoChan/dineflow/payment/app/command"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

type Consumer struct {
	app app.Application
}

func NewConsumer(app app.Application) *Consumer {
	return &Consumer{app: app}
}

func (c *Consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
	if err != nil {
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
	tr := otel.Tracer("rabbit-mq")
	_, span := tr.Start(ctx, fmt.Sprintf("rabbit-mq.%s.consume", q.Name))
	defer span.End()

	logging.Infof(ctx, nil, "Payment received a message from %s, msg=%v", q.Name, string(msg.Body))

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

	// TODO: ACL Cleaning
	order := &orderpb.Order{}
	if err = json.Unmarshal(msg.Body, order); err != nil {
		err = errors.Wrap(err, "failed to unmarshal msg to order")
		return
	}

	if _, err = c.app.Commands.CreatePayment.Handle(ctx, command.CreatePayment{
		Order: order,
	}); err != nil {
		err = errors.Wrap(err, "failed to create payment")
		if err = broker.HandleRetry(ctx, ch, &msg); err != nil {
			err = errors.Wrapf(err, "retry_error handling retry, msgID=%s, err=%v", msg.MessageId, err)
		}
		return
	}

	span.AddEvent("payment.created")
}
