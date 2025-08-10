package command

import (
	"context"

	"github.com/FacundoChan/dineflow/common/decorator"
	"github.com/FacundoChan/dineflow/common/logging"
	domain "github.com/FacundoChan/dineflow/order/domain/order"
	"github.com/sirupsen/logrus"
)

type UpdateOrder struct {
	Order      *domain.Order                                               `json:"order"`
	UpdateFunc func(context.Context, *domain.Order) (*domain.Order, error) `json:"-"`
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
	var err error
	defer logging.WhenCommandExecute(ctx, "UpdateOrderHandler", cmd, err)

	if cmd.UpdateFunc == nil {
		logrus.Panicf("updateOrderHandler called with nil UpdateFunc, orderID=%#v", cmd.Order)
		// cmd.UpdateFunc = func(_ context.Context, order *domain.Order) (*domain.Order, error) {
		// 	return order, nil // do nothing
		// }
	}
	if err = c.orderRepo.Update(ctx, cmd.Order, cmd.UpdateFunc); err != nil {
		return nil, err
	}
	return nil, nil
}
