package command

import (
	"context"

	"github.com/FacundoChan/dineflow/common/decorator"
	"github.com/FacundoChan/dineflow/common/genproto/orderpb"
	"github.com/FacundoChan/dineflow/common/logging"
	"github.com/FacundoChan/dineflow/common/tracing"
	"github.com/FacundoChan/dineflow/payment/domain"
	"github.com/sirupsen/logrus"
)

// TODO: ACL Cleaning
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
	var err error
	defer logging.WhenCommandExecute(ctx, "CreatePaymentHandler", cmd, err)

	_, span := tracing.Start(ctx, "stripe_processor.create_payment_link")
	defer span.End()

	link, err := c.processor.CreatePaymentLink(ctx, cmd.Order)
	if err != nil {
		return "", err
	}

	newOrder := &orderpb.Order{
		ID:          cmd.Order.ID,
		CustomerID:  cmd.Order.CustomerID,
		Status:      "waiting_for_payment",
		Items:       cmd.Order.Items,
		PaymentLink: link,
	}

	err = c.orderGRPC.UpdateOrder(ctx, newOrder)
	if err != nil {
		logrus.WithError(err).Debug("UpdateOrder failed")
	}

	return link, err
}
