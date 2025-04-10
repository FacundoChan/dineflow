package grpc

import (
	"context"

	"github.com/FacundoChan/gorder-v1/common/genproto/orderpb"
	"github.com/FacundoChan/gorder-v1/common/genproto/stockpb"
	"github.com/sirupsen/logrus"
)

type StockGRPC struct {
	client stockpb.StockServiceClient
}

func NewStockGRPC(client stockpb.StockServiceClient) *StockGRPC {
	return &StockGRPC{client: client}
}

func (s StockGRPC) GetItems(ctx context.Context, itemsID []string) ([]*orderpb.Item, error) {
	response, err := s.client.GetItems(ctx, &stockpb.GetItemsRequest{
		ItemIDs: itemsID,
	})
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	return response.Items, nil

}

func (s StockGRPC) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error) {
	response, err := s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{
		Items: items,
	})
	logrus.Info("gRPC CheckIfItemsInStock response:", response)
	return response, err
}

func (s StockGRPC) GetAllItems(ctx context.Context) (*stockpb.GetAllItemsResponse, error) {
	logrus.Debugf("gRPC GetAllItems")
	response, err := s.client.GetAllItems(ctx, &stockpb.GetAllItemsRequest{})
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	logrus.Info("gRPC GetAllItems response:", response)
	return response, err
}
