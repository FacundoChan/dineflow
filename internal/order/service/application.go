package service

import (
	"context"
	"github.com/FacundoChan/gorder-v1/common/metrics"
	"github.com/FacundoChan/gorder-v1/order/adapters"
	"github.com/FacundoChan/gorder-v1/order/app"
	"github.com/FacundoChan/gorder-v1/order/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.TodoMetrics{}
	return app.Application{
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricsClient),
		},
	}
}
