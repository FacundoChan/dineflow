package query

import (
	"context"

	"github.com/FacundoChan/gorder-v1/common/decorator"
	domain "github.com/FacundoChan/gorder-v1/stock/domain/stock"
	"github.com/FacundoChan/gorder-v1/stock/entity"
	"github.com/sirupsen/logrus"
)

type GetItems struct {
	ItemIDs []string
}

type GetItemsHandler decorator.QueryHandler[GetItems, []*entity.Item]

type getItemsHandler struct {
	stockRepo domain.Repository
}

func NewGetItemsHandler(stockRepo domain.Repository, logger *logrus.Entry, metricsClient decorator.MetricsClient) GetItemsHandler {
	if stockRepo == nil {
		panic("stockRepo is nil")
	}
	return decorator.ApplyQueryDecorators[GetItems, []*entity.Item](
		getItemsHandler{stockRepo: stockRepo},
		logger,
		metricsClient,
	)
}

func (g getItemsHandler) Handle(ctx context.Context, query GetItems) ([]*entity.Item, error) {
	items, err := g.stockRepo.GetItems(ctx, query.ItemIDs)
	if err != nil {
		logrus.Debugf("get items err: %v", err)
		return nil, err
	}
	return items, nil
}
