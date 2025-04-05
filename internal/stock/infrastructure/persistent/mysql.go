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
	ID        int64     `json:"column:id"`
	ProductID string    `json:"column:product_id"`
	Quantity  int64     `json:"column:quantity"`
	CreatedAt time.Time `json:"column:created_at"`
	UpdatedAt time.Time `json:"column:updated_at"`
}

func (m StockModel) TableName() string {
	return "order_stock"
}

func (d MySQL) BatchGetStockByProductIDs(ctx context.Context, productIDs []string) ([]StockModel, error) {
	var result []StockModel
	tx := d.db.WithContext(ctx).Where("product_id IN ?", productIDs).Find(&result)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return result, nil
}
