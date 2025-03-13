package command

import (
	"context"
	"github.com/FacundoChan/gorder-v1/common/genproto/orderpb"
)

type OrderService interface {
	UpdateOrder(ctx context.Context, order *orderpb.Order) error
}
