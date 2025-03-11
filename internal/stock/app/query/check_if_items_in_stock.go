package query

import (
	"context"
	"github.com/FacundoChan/gorder-v1/common/decorator"
	"github.com/FacundoChan/gorder-v1/common/genproto/orderpb"
	domain "github.com/FacundoChan/gorder-v1/stock/domain/stock"
	"github.com/sirupsen/logrus"
)

type CheckIfItemsInStock struct {
	Items []*orderpb.ItemWithQuantity
}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*orderpb.Item]

type checkIfItemsInStockHandler struct {
	stockRepo domain.Repository
}

func NewCheckIfItemsInStockHandler(stockRepo domain.Repository, logger *logrus.Entry, metricsClient decorator.MetricsClient) CheckIfItemsInStockHandler {
	if stockRepo == nil {
		panic("nil orderRepo")
	}
	return decorator.ApplyQueryDecorators[CheckIfItemsInStock, []*orderpb.Item](
		checkIfItemsInStockHandler{stockRepo: stockRepo},
		logger,
		metricsClient,
	)
}

func (c checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*orderpb.Item, error) {
	var res []*orderpb.Item
	for _, item := range query.Items {
		res = append(res, &orderpb.Item{
			ID:       item.ID,
			Quantity: item.Quantity,
		})
	}
	return res, nil
}
