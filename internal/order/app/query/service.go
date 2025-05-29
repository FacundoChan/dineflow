package query

import (
	"context"

	"github.com/FacundoChan/dineflow/common/genproto/orderpb"
	"github.com/FacundoChan/dineflow/common/genproto/stockpb"
)

type StockService interface {
	// rpc GetItems(GetItemsRequest) returns (GetItemsResponse);
	GetItems(ctx context.Context, itemsID []string) ([]*orderpb.Item, error)
	// rpc CheckIfItemsInStock(CheckIfItemsInStockRequest) returns (CheckIfItemsInStockResponse);
	CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error)
	// rpc GetAllProducts(GetAllProductsRequest) returns (GetAllProductsResponse);
	GetAllProducts(ctx context.Context) (*stockpb.GetAllProductsResponse, error)
}
