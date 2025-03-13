package consumer

import (
	"context"
	"encoding/json"

	"github.com/FacundoChan/gorder-v1/common/broker"
	"github.com/FacundoChan/gorder-v1/common/genproto/orderpb"
	"github.com/FacundoChan/gorder-v1/payment/app"
	"github.com/FacundoChan/gorder-v1/payment/app/command"
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
	q, err := ch.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Warnf("Failed to consume message: queue=%s, err=%s", q.Name, err)
	}

	var forever chan struct{}

	go func() {
		for msg := range msgs {
			c.handleMessage(msg, q, ch)
		}
		<-forever
	}()

}

func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue, ch *amqp.Channel) {
	logrus.Infof("Received message from %s: %s", q.Name, msg.Body)

	order := &orderpb.Order{}
	if err := json.Unmarshal(msg.Body, order); err != nil {
		logrus.Infof("failed to unmarshal msg to order, err: %v", err)
		_ = msg.Nack(false, false)
		return
	}

	if _, err := c.app.Commands.CreatePayment.Handle(context.TODO(), command.CreatePayment{
		Order: order,
	}); err != nil {
		// TODO: Retry
		logrus.Infof("failed to create order, err: %v", err)
		_ = msg.Nack(false, false)
		return
	}

	_ = msg.Ack(false)
	logrus.Info("successfully consumed")
}
