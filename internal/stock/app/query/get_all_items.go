package query

import (
	"context"
	"github.com/FacundoChan/gorder-v1/common/decorator"
	"github.com/FacundoChan/gorder-v1/common/utils"
	domain "github.com/FacundoChan/gorder-v1/stock/domain/stock"
	"github.com/FacundoChan/gorder-v1/stock/entity"
	"github.com/sirupsen/logrus"
)

type GetAllItems struct {
}

type GetAllItemsHandler decorator.QueryHandler[GetAllItems, []*entity.Item]

type getAllItemsHandler struct {
	stockRepo domain.Repository
}

func NewGetAllItemsHandler(stockRepo domain.Repository, logger *logrus.Entry, metricsClient decorator.MetricsClient) GetAllItemsHandler {
	if stockRepo == nil {
		panic("stockRepo is nil")
	}
	return decorator.ApplyQueryDecorators[GetAllItems, []*entity.Item](
		getAllItemsHandler{stockRepo: stockRepo},
		logger,
		metricsClient,
	)
}

// Handle implements decorator.QueryHandler.
func (g getAllItemsHandler) Handle(ctx context.Context, query GetAllItems) ([]*entity.Item, error) {
	allItems, err := g.stockRepo.GetAllItems(ctx)
	logrus.WithFields(logrus.Fields{
		"items": utils.ToString(allItems),
	}).Debug("[getAllItemsHandler.Handle]")
	if err != nil {
		logrus.Debugf("get all items err: %v", err)
		return nil, err
	}
	return allItems, nil
}
