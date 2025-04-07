package adapters

import (
	"context"
	"fmt"

	"github.com/FacundoChan/gorder-v1/common/utils"
	"github.com/FacundoChan/gorder-v1/stock/entity"
	"github.com/FacundoChan/gorder-v1/stock/infrastructure/persistent"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MySQLStockRepository struct {
	db *persistent.MySQL
}

func NewMySQLStockRepository(db *persistent.MySQL) *MySQLStockRepository {
	return &MySQLStockRepository{db: db}
}

func (m MySQLStockRepository) GetItems(ctx context.Context, ids []string) ([]*entity.Item, error) {
	var (
		result []*entity.Item
	)

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

func (m MySQLStockRepository) GetStock(ctx context.Context, ids []string) ([]*entity.ItemWithQuantity, error) {
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
	return m.db.StartTransaction(func(tx *gorm.DB) error {
		defer func() {
			if err != nil {
				logrus.Warnf("update transaction error: %v", err)
			}
		}()

		var dest []*persistent.StockModel
		// HACK: table name should be variable
		if err = tx.Table("order_stock").Where("product_id IN ?", getIDFromEntities(data)).Find(&dest).Error; err != nil {
			return errors.Wrap(err, "failed to get product_id")
		}
		existing := m.unmarshalFromDatabase(dest)
		logrus.WithFields(logrus.Fields{
			"existing": utils.ToString(existing),
		}).Debug("[existing]")

		updated, err := updateFunc(ctx, existing, data)
		if err != nil {
			return err
		}
		logrus.WithFields(logrus.Fields{
			"updated": utils.ToString(updated),
		}).Debug("[updated]")

		for _, updatedData := range updated {
			// HACK: table name should be variable
			if err = tx.Table("order_stock").Where("product_id = ?", updatedData.ID).Update("quantity", updatedData.Quantity).Error; err != nil {
				return errors.Wrap(err, fmt.Sprintf("unable to update %v+", updatedData))
			}
		}

		return nil
	})

}

func getIDFromEntities(data []*entity.ItemWithQuantity) []string {
	var result []string
	for _, d := range data {
		result = append(result, d.ID)
	}
	return result
}

func (m MySQLStockRepository) unmarshalFromDatabase(dest []*persistent.StockModel) []*entity.ItemWithQuantity {
	var result []*entity.ItemWithQuantity
	for _, i := range dest {
		result = append(result, &entity.ItemWithQuantity{
			ID:       i.ProductID,
			Quantity: int32(i.Quantity),
		})
	}
	return result
}
