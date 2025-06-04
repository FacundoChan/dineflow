package grpc

import (
	"context"

	"github.com/FacundoChan/dineflow/common/genproto/orderpb"
	"github.com/FacundoChan/dineflow/common/genproto/stockpb"
	"github.com/FacundoChan/dineflow/common/logging"
	"github.com/sirupsen/logrus"
)

type StockGRPC struct {
	client stockpb.StockServiceClient
}

func NewStockGRPC(client stockpb.StockServiceClient) *StockGRPC {
	return &StockGRPC{client: client}
}

func (s StockGRPC) GetItems(ctx context.Context, itemsID []string) (items []*orderpb.Item, err error) {
	_, dLog := logging.WhenRequest(ctx, "StockGRPC.GetItems", items)
	defer dLog(items, &err)

	response, err := s.client.GetItems(ctx, &stockpb.GetItemsRequest{
		ItemIDs: itemsID,
	})
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	return response.Items, nil

}

func (s StockGRPC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (response *stockpb.CheckIfItemsInStockResponse, err error) {
	_, dLog := logging.WhenRequest(ctx, "StockGRPC.CheckIfItemsInStock", items)
	defer dLog(response, &err)
	return s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{
		Items: items,
	})
}

func (s StockGRPC) GetAllProducts(ctx context.Context) (*stockpb.GetAllProductsResponse, error) {
	logrus.Debugf("gRPC GetAllProducts")
	response, err := s.client.GetAllProducts(ctx, &stockpb.GetAllProductsRequest{})
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	logrus.Info("gRPC GetAllProducts response:", response)
	return response, err
}
