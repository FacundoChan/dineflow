package service

import (
	"context"

	"github.com/FacundoChan/dineflow/common/metrics"
	"github.com/FacundoChan/dineflow/stock/adapters"
	"github.com/FacundoChan/dineflow/stock/app"
	"github.com/FacundoChan/dineflow/stock/app/query"
	"github.com/FacundoChan/dineflow/stock/infrastructure/integration"
	"github.com/FacundoChan/dineflow/stock/infrastructure/persistent"
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
	metricsClient := metrics.NewPrometheusMetricsClient(&metrics.PrometheusMetricsClientConfig{
		ServiceName: viper.GetString("stock.service-name"),
		Host:        viper.GetString("stock.metrics_export_addr"),
	})
	logrus.Info("PrometheusMetricsClient initialized")

	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, stripAPI, logger, metricsClient),
			GetItems:            query.NewGetItemsHandler(stockRepo, logger, metricsClient),
			GetAllProducts:      query.NewGetAllProductsHandler(stockRepo, logger, metricsClient),
		},
	}
}
