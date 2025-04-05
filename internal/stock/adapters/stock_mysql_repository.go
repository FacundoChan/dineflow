package adapters

import (
	"context"

	"github.com/FacundoChan/gorder-v1/stock/entity"
	"github.com/FacundoChan/gorder-v1/stock/infrastructure/persistent"
)

type MySQLStockRepository struct {
	db *persistent.MySQL
}

func NewMySQLStockRepository(db *persistent.MySQL) *MySQLStockRepository {
	return &MySQLStockRepository{db: db}
}

func (m MySQLStockRepository) GetItems(ctx context.Context, ids []string) ([]*entity.Item, error) {
	//TODO implement me
	panic("implement me")
}

func (m MySQLStockRepository) GetStock(ctx context.Context, ids []string) ([]*entity.ItemWithQuantity, error) {
	data, err := m.db.BatchGetStockByProductIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	var result []*entity.ItemWithQuantity
	for _, d := range data {
		result = append(result, &entity.ItemWithQuantity{
			ID:       d.ProductID,
			Quantity: int32(d.Quantity),
		})
	}
	return result, nil
}
