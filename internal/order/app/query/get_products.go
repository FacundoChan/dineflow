package query

import (
	"context"

	"github.com/FacundoChan/dineflow/common/decorator"
	"github.com/FacundoChan/dineflow/common/format"
	"github.com/FacundoChan/dineflow/order/convertor"
	"github.com/FacundoChan/dineflow/order/entity"
	"github.com/sirupsen/logrus"
)

type GetProducts struct {
}

type GetProductsResult struct {
	Products []*entity.Product
}

type GetProductsHandler decorator.QueryHandler[GetProducts, *GetProductsResult]

type getProductsHandler struct {
	stockGRPC StockService
}

func NewGetProductsHandler(
	stockGRPC StockService,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient) GetProductsHandler {

	if stockGRPC == nil {
		logrus.Panic("stockGRPC is nil")
	}
	return decorator.ApplyQueryDecorators[GetProducts, *GetProductsResult](
		getProductsHandler{
			stockGRPC: stockGRPC,
		},
		logger,
		metricsClient,
	)
}

func (g getProductsHandler) Handle(ctx context.Context, query GetProducts) (*GetProductsResult, error) {
	response, err := g.stockGRPC.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}
	result := convertor.NewProductConvertor().ProtosToEntities(response.Products)

	logrus.WithFields(logrus.Fields{
		"result": format.ToString(result),
	}).Debug("[getProductsHandler.Handle]")

	return &GetProductsResult{Products: result}, nil
}
