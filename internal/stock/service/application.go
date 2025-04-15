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
	"github.com/spf13/viper"
)

func NewApplication(ctx context.Context) app.Application {
	// stockRepo := adapters.NewMemoryStockRepository()
	db := persistent.NewMySQL()
	stockRepo := adapters.NewMySQLStockRepository(db)
	stripeKey := viper.GetString("stripe-key")
	if stripeKey == "" {
		logrus.Fatal("stripe-key is empty")
	} else {
		logrus.Info("[NewStripeAPI] stripe-key found.")
	}

	stripAPI := integration.NewStripeAPI(stripeKey)
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.TodoMetrics{}

	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, stripAPI, logger, metricsClient),
			GetItems:            query.NewGetItemsHandler(stockRepo, logger, metricsClient),
			GetAllProducts:      query.NewGetAllProductsHandler(stockRepo, logger, metricsClient),
		},
	}
}
