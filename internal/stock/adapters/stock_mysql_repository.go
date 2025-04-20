package adapters

import (
	"context"

	"github.com/FacundoChan/gorder-v1/common/utils"
	domain "github.com/FacundoChan/gorder-v1/stock/domain/stock"
	"github.com/FacundoChan/gorder-v1/stock/entity"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type MySQLStockRepository struct {
	db domain.StockDB
}

func NewMySQLStockRepository(db domain.StockDB) *MySQLStockRepository {
	return &MySQLStockRepository{db: db}
}

func (m MySQLStockRepository) GetItems(ctx context.Context, ids []string) ([]*entity.Item, error) {
	var (
		result []*entity.Item
	)

	// query := builder.NewStock().ProductIDs(ids...)
	// data, err := m.db.BatchGetStockByProductIDs(ctx, query)

	data, err := m.db.BatchGetStockByProductIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	for _, d := range data {
		result = append(result, &entity.Item{
			ID:       d.ProductID,
			Quantity: int32(d.Quantity),
		})
	}

	return result, nil
}

func (m MySQLStockRepository) GetAllProducts(ctx context.Context) ([]*entity.Product, error) {
	var (
		result []*entity.Product
	)

	data, err := m.db.GetAllStockProducts(ctx)
	if err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"items": utils.ToString(data),
	}).Debug("[data]")

	for _, d := range data {
		result = append(result, &entity.Product{
			ID:       d.ProductID,
			Name:     d.Name,
			Quantity: int32(d.Quantity),
			Price:    d.Price,
			ImgUrls:  d.ImgUrls,
		})
	}

	return result, nil
}

func (m MySQLStockRepository) GetStock(ctx context.Context, ids []string) ([]*entity.ItemWithQuantity, error) {
	// query := builder.NewStock().ProductIDs(ids...)
	// data, err := m.db.BatchGetStockByProductIDs(ctx, query)
	data, err := m.db.BatchGetStockByProductIDs(ctx, ids)
	if err != nil {
		return nil, errors.Wrap(err, "BatchGetStockByProductIDs failed")
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

func (m MySQLStockRepository) UpdateStock(ctx context.Context,
	data []*entity.ItemWithQuantity,
	updateFunc func(context.Context, []*entity.ItemWithQuantity, []*entity.ItemWithQuantity) ([]*entity.ItemWithQuantity, error)) (err error) {
	logrus.Debug("mysql_repo_UpdateStock called")
	return m.db.UpdateStockTransaction(ctx, data, updateFunc)
}

// func getIDFromEntities(data []*entity.ItemWithQuantity) []string {
// 	var result []string
// 	for _, d := range data {
// 		result = append(result, d.ID)
// 	}
// 	return result
// }

// func (m MySQLStockRepository) unmarshalFromDatabase(dest []*persistent.StockModel) []*entity.ItemWithQuantity {
// 	var result []*entity.ItemWithQuantity
// 	for _, i := range dest {
// 		result = append(result, &entity.ItemWithQuantity{
// 			ID:       i.ProductID,
// 			Quantity: int32(i.Quantity),
// 		})
// 	}
// 	return result
// }
