package persistent

import (
	"context"
	"fmt"
	"time"

	"github.com/FacundoChan/dineflow/common/format"
	"github.com/FacundoChan/dineflow/common/logging"
	"github.com/FacundoChan/dineflow/stock/entity"
	"github.com/FacundoChan/dineflow/stock/infrastructure/persistent/builder"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MySQL struct {
	db *gorm.DB
}

func NewMySQL() *MySQL {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		viper.GetString("mysql.db-name"),
	)
	logrus.Debugf("[MYSQL] dsn = %s", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Panicf("failed to connect mysql, err=%v", err)
	}
	return &MySQL{db: db}
}

func NewMySQLWithDB(db *gorm.DB) *MySQL {
	return &MySQL{db: db}
}

type StockModel struct {
	ID          int64     `gorm:"column:id"`
	ProductID   string    `gorm:"column:product_id"`
	Name        string    `gorm:"column:name"`
	Quantity    int64     `gorm:"column:quantity"`
	Price       float32   `gorm:"column:price"`
	Description string    `gorm:"column:description"`
	Version     int64     `gorm:"column:version"`
	ImgUrls     []string  `gorm:"-"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (m StockModel) TableName() string {
	return "order_stock"
}

func (d *MySQL) StartTransaction(fc func(tx *gorm.DB) error) error {
	return d.db.Transaction(fc)
}

func (d *MySQL) UpdateStockTransaction(ctx context.Context, data []*entity.ItemWithQuantity,
	updateFunc func(c context.Context, existing []*entity.ItemWithQuantity, query []*entity.ItemWithQuantity) ([]*entity.ItemWithQuantity, error)) error {
	return d.StartTransaction(func(tx *gorm.DB) (err error) {

		defer func() {
			if err != nil {
				logrus.Warnf("update transaction error: %v", err)
			}
		}()

		err = d.UpdateStockPessimistic(ctx, tx, data, updateFunc)
		// err = d.UpdateStockOptimistic(ctx, tx, data, updateFunc)

		return err
	})

}

func (d *MySQL) UpdateStockOptimistic(
	ctx context.Context,
	tx *gorm.DB,
	data []*entity.ItemWithQuantity,
	updateFunc func(c context.Context,
		existing []*entity.ItemWithQuantity,
		query []*entity.ItemWithQuantity,
	) ([]*entity.ItemWithQuantity, error)) error {
	var dest []*StockModel
	// HACK: table name should be variable
	if err := builder.NewStock().ProductIDs(getIDFromEntities(data)...).
		Fill(tx.Table("order_stock")).Find(&dest).Error; err != nil {
		return errors.Wrap(err, "failed to get product_id with lock")
	}

	existing := d.unmarshalFromDatabase(dest)
	logrus.WithFields(logrus.Fields{
		"existing": format.ToString(existing),
	}).Debug("[existing]")

	updated, err := updateFunc(ctx, existing, data)
	if err != nil {
		return err
	}

	for _, queryData := range data {
		// var newestRecord entity.StockModel
		var newestRecord StockModel

		// HACK: table name should be variable
		if err := builder.NewStock().ProductIDs(queryData.ID).
			Fill(tx.Table("order_stock")).First(&newestRecord).Error; err != nil {
			return err
		}

		if err := builder.NewStock().ProductIDs(queryData.ID).Versions(newestRecord.Version).QuantityGreaterEqual(queryData.Quantity).
			Fill(tx.Table("order_stock")).
			Updates(map[string]any{
				"quantity": gorm.Expr("quantity - ?", queryData.Quantity),
				"version":  newestRecord.Version + 1,
			}).Error; err != nil {
			return err
		}
	}

	logrus.WithFields(logrus.Fields{
		"updated": format.ToString(updated),
	}).Debug("[updated]")

	return nil
}

func (d *MySQL) UpdateStockPessimistic(
	ctx context.Context,
	tx *gorm.DB,
	data []*entity.ItemWithQuantity,
	updateFunc func(c context.Context,
		existing []*entity.ItemWithQuantity,
		query []*entity.ItemWithQuantity,
	) ([]*entity.ItemWithQuantity, error)) error {
	_, deferLog := logging.WhenMySQL(ctx, "UpdateStockPessimistic", data)
	var updateResult []*entity.ItemWithQuantity
	var err error

	defer func() {
		deferLog(updateResult, &err)
	}()

	var dest []*StockModel
	// HACK: table name should be variable
	if err = builder.NewStock().ProductIDs(getIDFromEntities(data)...).
		ForUpdate().
		Fill(tx.Table("order_stock")).
		Find(&dest).Error; err != nil {
		return errors.Wrap(err, "failed to get product_id with lock")
	}

	existing := d.unmarshalFromDatabase(dest)
	logrus.WithFields(logrus.Fields{
		"existing": format.ToString(existing),
	}).Debug("[existing]")

	updateResult, err = updateFunc(ctx, existing, data)
	if err != nil {
		return err
	}
	logrus.WithFields(logrus.Fields{
		"updated": format.ToString(updateResult),
	}).Debug("[updated]")

	for _, updatedData := range updateResult {
		for _, query := range data {
			if query.ID == updatedData.ID {
				// HACK: table name should be variable
				if err = builder.NewStock().ProductIDs(updatedData.ID).
					QuantityGreaterEqual(query.Quantity).
					Fill(tx.Table("order_stock")).
					Update("quantity", gorm.Expr("quantity - ?", query.Quantity)).Error; err != nil {
					return errors.Wrap(err, fmt.Sprintf("unable to update %v+", updatedData))
				}
			}
		}

	}

	return nil
}

func (d *MySQL) BatchGetStockByProductIDs(ctx context.Context, productIDs []string) ([]entity.StockModel, error) {
	var result []StockModel
	query := builder.NewStock()
	_, deferLog := logging.WhenMySQL(ctx, "BatchGetStockByProductIDs", query)

	tx := query.Fill(d.db.WithContext(ctx).Clauses(clause.Returning{})).Find(&result)

	defer deferLog(result, &tx.Error)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return d.PersistentsToEntities(result), nil
}

func (d *MySQL) GetStocksByPage(ctx context.Context, offset int, limit int) ([]entity.StockModel, error) {
	// For handler
	// offset := (page - 1) * pageSize
	// limit := pageSize
	// stocks, err := repo.GetStocksByPage(ctx, offset, limit)

	var result []StockModel

	tx := d.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&result)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return d.PersistentsToEntities(result), nil
}

func (d *MySQL) GetAllStockProducts(ctx context.Context) ([]entity.StockModel, error) {
	var result []StockModel

	// TODO: replaced with builder
	if err := d.db.WithContext(ctx).Limit(10).Find(&result).Error; err != nil {
		return nil, err
	}

	for i := range result {
		product := &result[i]
		var productImgsUrls []string
		if err := d.db.WithContext(ctx).Table("product_images").Where("product_id = ?", product.ProductID).Pluck("img_url", &productImgsUrls).Error; err != nil {
			return nil, err
		}

		product.ImgUrls = productImgsUrls
	}

	return d.PersistentsToEntities(result), nil
}

func (d *MySQL) Create(ctx context.Context, create *StockModel) error {
	_, deferLog := logging.WhenMySQL(ctx, "Create", create)
	var returning StockModel
	err := d.db.WithContext(ctx).Model(&returning).Clauses(clause.Returning{}).Create(create).Error
	defer deferLog(returning, &err)
	return err
}

func getIDFromEntities(data []*entity.ItemWithQuantity) []string {
	var result []string
	for _, d := range data {
		result = append(result, d.ID)
	}
	return result
}

func (d *MySQL) unmarshalFromDatabase(dest []*StockModel) []*entity.ItemWithQuantity {
	var result []*entity.ItemWithQuantity
	for _, i := range dest {
		result = append(result, &entity.ItemWithQuantity{
			ID:       i.ProductID,
			Quantity: int32(i.Quantity),
		})
	}
	return result
}

func (d *MySQL) PersistentsToEntities(s []StockModel) []entity.StockModel {
	var res []entity.StockModel
	for _, p := range s {
		res = append(res, d.PersistentToEntity(p))
	}
	return res
}

func (d *MySQL) PersistentToEntity(p StockModel) entity.StockModel {
	return entity.StockModel{
		ID:          p.ID,
		ProductID:   p.ProductID,
		Name:        p.Name,
		Quantity:    p.Quantity,
		Price:       p.Price,
		Description: p.Description,
		Version:     p.Version,
		ImgUrls:     p.ImgUrls,
	}
}
