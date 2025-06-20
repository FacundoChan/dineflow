package adapters

import (
	"context"

	"github.com/FacundoChan/dineflow/common/genproto/orderpb"
	"github.com/FacundoChan/dineflow/common/tracing"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
)

type OrderRPC struct {
	client orderpb.OrderServiceClient
}

func NewOrderRPC(client orderpb.OrderServiceClient) *OrderRPC {
	return &OrderRPC{
		client: client,
	}
}

func (o OrderRPC) UpdateOrder(ctx context.Context, order *orderpb.Order) error {
	ctx, span := tracing.Start(ctx, "order_grpc.update_order")
	defer span.End()

	_, err := o.client.UpdateOrder(ctx, order)
	if err != nil {
		logrus.Errorf("payment adapter failed to update order: %v", err)
		return status.Convert(err).Err()
	}
	return nil
}
