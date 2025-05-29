package query

import (
	"context"

	"github.com/FacundoChan/dineflow/common/decorator"
	"github.com/FacundoChan/dineflow/common/utils"
	domain "github.com/FacundoChan/dineflow/stock/domain/stock"
	"github.com/FacundoChan/dineflow/stock/entity"
	"github.com/sirupsen/logrus"
)

type GetAllProducts struct {
}

type GetAllProductsHandler decorator.QueryHandler[GetAllProducts, []*entity.Product]

type getAllProductsHandler struct {
	stockRepo domain.Repository
}

func NewGetAllProductsHandler(stockRepo domain.Repository, logger *logrus.Entry, metricsClient decorator.MetricsClient) GetAllProductsHandler {
	if stockRepo == nil {
		panic("stockRepo is nil")
	}
	return decorator.ApplyQueryDecorators[GetAllProducts, []*entity.Product](
		getAllProductsHandler{stockRepo: stockRepo},
		logger,
		metricsClient,
	)
}

// Handle implements decorator.QueryHandler.
func (g getAllProductsHandler) Handle(ctx context.Context, query GetAllProducts) ([]*entity.Product, error) {
	allProducts, err := g.stockRepo.GetAllProducts(ctx)
	logrus.WithFields(logrus.Fields{
		"products": utils.ToString(allProducts),
	}).Debug("[getAllItemsHandler.Handle]")
	if err != nil {
		logrus.Debugf("get all products err: %v", err)
		return nil, err
	}
	return allProducts, nil
}
