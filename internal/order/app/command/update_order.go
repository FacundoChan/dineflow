package command

import (
	"context"

	"github.com/FacundoChan/dineflow/common/decorator"
	domain "github.com/FacundoChan/dineflow/order/domain/order"
	"github.com/sirupsen/logrus"
)

type UpdateOrder struct {
	Order      *domain.Order
	UpdateFunc func(context.Context, *domain.Order) (*domain.Order, error)
}

type UpdateOrderHandler decorator.CommandHandler[UpdateOrder, interface{}]

type updateOrderHandler struct {
	orderRepo domain.Repository
	//stockGRPC
}

func NewUpdateOrderHandler(orderRepo domain.Repository, logger *logrus.Entry, metricsClient decorator.MetricsClient) UpdateOrderHandler {
	if orderRepo == nil {
		panic("orderRepo is nil")
	}
	return decorator.ApplyCommandDecorators[UpdateOrder, interface{}](
		updateOrderHandler{orderRepo: orderRepo},
		logger,
		metricsClient,
	)
}

func (c updateOrderHandler) Handle(ctx context.Context, cmd UpdateOrder) (interface{}, error) {
	if cmd.UpdateFunc == nil {
		logrus.Warnf("updateOrderHandler called with nil UpdateFunc, orderID=%#v", cmd.Order)
		cmd.UpdateFunc = func(_ context.Context, order *domain.Order) (*domain.Order, error) {
			return order, nil // do nothing
		}
	}
	err := c.orderRepo.Update(ctx, cmd.Order, cmd.UpdateFunc)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
