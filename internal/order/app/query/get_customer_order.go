package query

import (
	"context"

	"github.com/FacundoChan/gorder-v1/common/decorator"
	domain "github.com/FacundoChan/gorder-v1/order/domain/order"
	"github.com/sirupsen/logrus"
)

type GetCustomerOrder struct {
	CustomerID string
	OrderID    string
}

type GetCustomerOrderHandler decorator.QueryHandler[GetCustomerOrder, *domain.Order]

type getCustomerOrderHandler struct {
	orderRepo domain.Repository
}

func NewGetCustomerOrderHandler(orderRepo domain.Repository, logger *logrus.Entry, metricsClient decorator.MetricsClient) GetCustomerOrderHandler {
	if orderRepo == nil {
		panic("nil orderRepo")
	}
	return decorator.ApplyQueryDecorators[GetCustomerOrder, *domain.Order](
		getCustomerOrderHandler{orderRepo: orderRepo},
		logger,
		metricsClient,
	)
}

func (g getCustomerOrderHandler) Handle(ctx context.Context, query GetCustomerOrder) (*domain.Order, error) {
	o, err := g.orderRepo.Get(ctx, query.OrderID, query.CustomerID)
	if err != nil {
		return nil, err
	}
	return o, nil
}
