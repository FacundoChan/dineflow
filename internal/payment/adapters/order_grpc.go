package adapters

import (
	"context"

	"github.com/FacundoChan/gorder-v1/common/genproto/orderpb"
	"github.com/FacundoChan/gorder-v1/common/tracing"
	"github.com/sirupsen/logrus"
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
		logrus.Infof("payment adapter failed to update order: %v", err)
	}
	return nil
}
