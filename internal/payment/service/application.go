package service

import (
	"context"

	grpcClient "github.com/FacundoChan/gorder-v1/common/client"
	"github.com/FacundoChan/gorder-v1/common/metrics"
	"github.com/FacundoChan/gorder-v1/payment/adapters"
	"github.com/FacundoChan/gorder-v1/payment/app"
	"github.com/FacundoChan/gorder-v1/payment/app/command"
	"github.com/FacundoChan/gorder-v1/payment/domain"
	"github.com/FacundoChan/gorder-v1/payment/infrastructure/processor"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	orderClient, closeOrderClient, err := grpcClient.NewOrderGRPCClient(ctx)
	if err != nil {
		logrus.Panic(err)
	}
	orderGRPC := adapters.NewOrderRPC(orderClient)
	stripeProcessor := processor.NewStripeProcessor(viper.GetString("stripe-key"))

	return newApplication(ctx, orderGRPC, stripeProcessor), func() {
		_ = closeOrderClient()
	}
}

func newApplication(_ context.Context, orderGRPC command.OrderService, processor domain.Processor) app.Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.TodoMetrics{}

	return app.Application{
		Commands: app.Commands{
			CreatePayment: command.NewCreatePaymentHandler(
				processor,
				orderGRPC,
				logger,
				metricsClient,
			),
		},
	}
}
