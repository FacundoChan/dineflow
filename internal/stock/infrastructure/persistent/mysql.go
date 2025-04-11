package persistent

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

type StockModel struct {
	ID          int64     `gorm:"column:id"`
	ProductID   string    `gorm:"column:product_id"`
	Name        string    `gorm:"column:name"`
	Quantity    int64     `gorm:"column:quantity"`
	Price       float32   `gorm:"column:price"`
	Description string    `gorm:"column:description"`
	ImgUrls     []string  `gorm:"-"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (m StockModel) TableName() string {
	return "order_stock"
}

func (d MySQL) StartTransaction(fc func(tx *gorm.DB) error) error {
	return d.db.Transaction(fc)
}

func (d MySQL) BatchGetStockByProductIDs(ctx context.Context, productIDs []string) ([]StockModel, error) {
	var result []StockModel
	tx := d.db.WithContext(ctx).Where("product_id IN ?", productIDs).Find(&result)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return result, nil
}

func (d MySQL) GetStocksByPage(ctx context.Context, offset int, limit int) ([]StockModel, error) {
	// For handler
	// offset := (page - 1) * pageSize
	// limit := pageSize
	// stocks, err := repo.GetStocksByPage(ctx, offset, limit)

	var result []StockModel

	tx := d.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&result)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return result, nil
}

func (d MySQL) GetAllStockProducts(ctx context.Context) ([]StockModel, error) {
	var result []StockModel

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

	return result, nil
}
