package service

import (
	"context"
	"github.com/FacundoChan/gorder-v1/common/metrics"
	"github.com/FacundoChan/gorder-v1/stock/adapters"
	"github.com/FacundoChan/gorder-v1/stock/app"
	"github.com/FacundoChan/gorder-v1/stock/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) app.Application {
	stockRepo := adapters.NewMemoryStockRepository()
	// Uncomment the next line to enable debug output - default is info and higher
	logrus.SetLevel(logrus.DebugLevel)
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.TodoMetrics{}

	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, logger, metricsClient),
			GetItems:            query.NewGetItemsHandler(stockRepo, logger, metricsClient),
		},
	}
}
