package service

import (
	"context"

	grpcClient "github.com/FacundoChan/dineflow/common/client"
	"github.com/FacundoChan/dineflow/common/metrics"
	"github.com/FacundoChan/dineflow/payment/adapters"
	"github.com/FacundoChan/dineflow/payment/app"
	"github.com/FacundoChan/dineflow/payment/app/command"
	"github.com/FacundoChan/dineflow/payment/domain"
	"github.com/FacundoChan/dineflow/payment/infrastructure/processor"
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
	metricsClient := metrics.NewPrometheusMetricsClient(&metrics.PrometheusMetricsClientConfig{
		ServiceName: viper.GetString("payment.service-name"),
		Host:        viper.GetString("payment.metrics_export_addr"),
	})
	logrus.Debugf("TOBUDELETED")

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
