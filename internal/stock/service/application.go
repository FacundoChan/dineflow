package service

import (
	"context"

	"github.com/FacundoChan/gorder-v1/common/metrics"
	"github.com/FacundoChan/gorder-v1/stock/adapters"
	"github.com/FacundoChan/gorder-v1/stock/app"
	"github.com/FacundoChan/gorder-v1/stock/app/query"
	"github.com/FacundoChan/gorder-v1/stock/infrastructure/integration"
	"github.com/FacundoChan/gorder-v1/stock/infrastructure/persistent"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) app.Application {
	// stockRepo := adapters.NewMemoryStockRepository()
	db := persistent.NewMySQL()
	stockRepo := adapters.NewMySQLStockRepository(db)
	stripAPI := integration.NewStripeAPI()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.TodoMetrics{}

	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, stripAPI, logger, metricsClient),
			GetItems:            query.NewGetItemsHandler(stockRepo, logger, metricsClient),
		},
	}
}
