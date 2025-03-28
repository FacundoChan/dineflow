package ports

import (
	"context"

	"github.com/FacundoChan/gorder-v1/common/genproto/orderpb"
	"github.com/FacundoChan/gorder-v1/order/app"
	"github.com/FacundoChan/gorder-v1/order/app/command"
	"github.com/FacundoChan/gorder-v1/order/app/query"
	"github.com/FacundoChan/gorder-v1/order/convertor"
	domain "github.com/FacundoChan/gorder-v1/order/domain/order"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
	_, err := G.app.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
		CustomerID: request.CustomerID,
		Items:      convertor.NewItemWithQuantityConvertor().ProtosToEntities(request.Items),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
	order, err := G.app.Queries.GetCustomerOrder.Handle(ctx, query.GetCustomerOrder{
		CustomerID: request.CustomerID,
		OrderID:    request.OrderID,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return convertor.NewOrderConvertor().EntityToProto(order), nil
}

func (G GRPCServer) UpdateOrder(ctx context.Context, request *orderpb.Order) (_ *emptypb.Empty, err error) {
	logrus.Infof("order_rpc request_in: %+v", request)
	order, err := domain.NewOrder(request.ID,
		request.CustomerID,
		request.Status,
		request.PaymentLink,
		convertor.NewItemConvertor().ProtosToEntities(request.Items),
	)
	if err != nil {
		err = status.Error(codes.Internal, err.Error())
		return nil, err
	}
	_, err = G.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
		Order: order,
		UpdateFunc: func(ctx context.Context, updatedOrder *domain.Order) (*domain.Order, error) {
			// TODO:
			return updatedOrder, nil
		},
	})
	return nil, err
}
