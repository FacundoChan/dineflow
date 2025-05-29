package ports

import (
	"context"

	"github.com/FacundoChan/dineflow/common/genproto/stockpb"
	"github.com/FacundoChan/dineflow/common/tracing"
	"github.com/FacundoChan/dineflow/common/utils"
	"github.com/FacundoChan/dineflow/stock/app"
	"github.com/FacundoChan/dineflow/stock/app/query"
	"github.com/FacundoChan/dineflow/stock/convertor"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
	_, span := tracing.Start(ctx, "GetItems")
	defer span.End()
	logrus.Info("rpc_request_in, stock.GetItems")
	defer func() {
		logrus.Info("rpc_request_out, stock.GetItems")
	}()

	items, err := G.app.Queries.GetItems.Handle(ctx, query.GetItems{
		ItemIDs: request.ItemIDs,
	})

	if err != nil {
		logrus.WithError(err).Error("rpc_request_err")
		return nil, err
	}

	return &stockpb.GetItemsResponse{
		Items: convertor.NewItemConvertor().EntitiesToProtos(items),
	}, nil
}

func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
	_, span := tracing.Start(ctx, "CheckIfItemsInStock")
	defer span.End()

	logrus.Info("rpc_request_in, stock.CheckIfItemsInStock")
	defer func() {
		logrus.Info("rpc_request_out, stock.CheckIfItemsInStock")
	}()
	items, err := G.app.Queries.CheckIfItemsInStock.Handle(ctx, query.CheckIfItemsInStock{
		Items: convertor.NewItemWithQuantityConvertor().ProtosToEntities(request.Items),
	})
	if err != nil {
		logrus.WithError(err).Error("rpc_request_err")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &stockpb.CheckIfItemsInStockResponse{
		InStock: 1,
		Items:   convertor.NewItemConvertor().EntitiesToProtos(items),
	}, nil
}

func (G GRPCServer) GetAllProducts(ctx context.Context, request *stockpb.GetAllProductsRequest) (*stockpb.GetAllProductsResponse, error) {
	var (
		err error
	)

	_, span := tracing.Start(ctx, "GetAllProducts")
	defer span.End()

	logrus.Info("rpc_request_in, stock.GetAllProducts")
	defer func() {
		logrus.Info("rpc_request_out, stock.GetAllProducts")
	}()

	allProducts, err := G.app.Queries.GetAllProducts.Handle(ctx, query.GetAllProducts{})
	logrus.WithFields(logrus.Fields{
		"products": utils.ToString(allProducts),
	}).Debug("[grpc.GetAllProducts]")
	if err != nil {
		logrus.WithError(err).Error("rpc_request_err")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &stockpb.GetAllProductsResponse{
		Products: convertor.NewProductConvertor().EntitiesToProtos(allProducts),
	}, nil
}
