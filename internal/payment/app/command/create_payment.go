package command

import (
	"context"

	"github.com/FacundoChan/gorder-v1/common/decorator"
	"github.com/FacundoChan/gorder-v1/common/genproto/orderpb"
	"github.com/FacundoChan/gorder-v1/common/tracing"
	"github.com/FacundoChan/gorder-v1/payment/domain"
	"github.com/sirupsen/logrus"
)

type CreatePayment struct {
	Order *orderpb.Order
}

type CreatePaymentHandler decorator.CommandHandler[CreatePayment, string]

type createPaymentHandler struct {
	processor domain.Processor
	orderGRPC OrderService
}

func NewCreatePaymentHandler(processor domain.Processor, orderGRPC OrderService, logger *logrus.Entry, metricsClient decorator.MetricsClient) CreatePaymentHandler {
	return decorator.ApplyCommandDecorators[CreatePayment, string](
		createPaymentHandler{
			processor: processor,
			orderGRPC: orderGRPC,
		},
		logger,
		metricsClient,
	)
}

// Handle implements decorator.CommandHandler.
func (c createPaymentHandler) Handle(ctx context.Context, cmd CreatePayment) (string, error) {
	_, span := tracing.Start(ctx, "stripe_processor.create_payment_link")
	defer span.End()

	link, err := c.processor.CreatePaymentLink(ctx, cmd.Order)
	if err != nil {
		return "", err
	}
	logrus.Infof("create payment link(%s) for order: %s", link, cmd.Order.CustomerID)
	newOrder := &orderpb.Order{
		ID:          cmd.Order.ID,
		CustomerID:  cmd.Order.CustomerID,
		Status:      "waiting_for_payment",
		Items:       cmd.Order.Items,
		PaymentLink: link,
	}

	logrus.Debug("UpdateOrder in createPaymentHandler")
	err = c.orderGRPC.UpdateOrder(ctx, newOrder)
	if err != nil {
		logrus.WithError(err).Debug("UpdateOrder failed")
	}

	return link, err
}
