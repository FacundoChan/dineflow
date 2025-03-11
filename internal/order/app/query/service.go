package query

import (
	"context"
	"github.com/FacundoChan/gorder-v1/common/genproto/orderpb"
	"github.com/FacundoChan/gorder-v1/common/genproto/stockpb"
)

type StockService interface {
	// rpc GetItems(GetItemsRequest) returns (GetItemsResponse);
	GetItems(ctx context.Context, itemsID []string) ([]*orderpb.Item, error)
	// rpc CheckIfItemsInStock(CheckIfItemsInStockRequest) returns (CheckIfItemsInStockResponse);
	CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error)
}
