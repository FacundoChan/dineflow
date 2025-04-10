package query

import (
	"context"

	"github.com/FacundoChan/gorder-v1/common/decorator"
	"github.com/FacundoChan/gorder-v1/order/convertor"
	"github.com/FacundoChan/gorder-v1/order/entity"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type GetProducts struct {
}

type GetProductsResult struct {
	Items []*entity.Item
}

type GetProductsHandler decorator.QueryHandler[GetProducts, *GetProductsResult]

type getProductsHandler struct {
	stockGRPC StockService
	channel   *amqp.Channel
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
	response, err := g.stockGRPC.GetAllItems(ctx)
	if err != nil {
		return nil, err
	}
	result := convertor.NewItemConvertor().ProtosToEntities(response.Items)
	return &GetProductsResult{Items: result}, nil
}
