package stock

import (
	"context"
	"fmt"
	"strings"

	"github.com/FacundoChan/dineflow/stock/entity"
)

type Repository interface {
	GetItems(ctx context.Context, ids []string) ([]*entity.Item, error)
	GetAllProducts(ctx context.Context) ([]*entity.Product, error)
	GetStock(ctx context.Context, ids []string) ([]*entity.ItemWithQuantity, error)
	UpdateStock(ctx context.Context, queryItems []*entity.ItemWithQuantity,
		fc func(ctx context.Context, existing []*entity.ItemWithQuantity, query []*entity.ItemWithQuantity) ([]*entity.ItemWithQuantity, error)) error
}

type StockDB interface {
	BatchGetStockByProductIDs(ctx context.Context, productIDs []string) ([]entity.StockModel, error)
	GetAllStockProducts(ctx context.Context) ([]entity.StockModel, error)
	UpdateStockTransaction(ctx context.Context, data []*entity.ItemWithQuantity,
		updateFunc func(c context.Context, existing []*entity.ItemWithQuantity, query []*entity.ItemWithQuantity) ([]*entity.ItemWithQuantity, error)) error
}

type NotFoundError struct {
	MissingIDs []string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("stock repository: not found: %s", strings.Join(e.MissingIDs, ","))
}

type ExceedStockError struct {
	NotEnoughItems []struct {
		ID   string
		Want int32
		Have int32
	}
}

func (e ExceedStockError) Error() string {
	var info []string
	for _, v := range e.NotEnoughItems {
		info = append(info, fmt.Sprintf("product_id=%s, want %v, have %v", v.ID, v.Want, v.Have))
	}
	return fmt.Sprintf("stock repository: not enough stock for [%s]", strings.Join(info, ", "))
}
