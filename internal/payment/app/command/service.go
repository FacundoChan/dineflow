package command

import (
	"context"

	"github.com/FacundoChan/dineflow/common/genproto/orderpb"
)

type OrderService interface {
	UpdateOrder(ctx context.Context, order *orderpb.Order) error
}
