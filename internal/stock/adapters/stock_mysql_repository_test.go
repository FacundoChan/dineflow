package adapters

import (
	"context"
	"fmt"
	"sync"
	"testing"

	_ "github.com/FacundoChan/dineflow/common/config"
	"github.com/FacundoChan/dineflow/stock/entity"
	"github.com/FacundoChan/dineflow/stock/infrastructure/persistent"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *persistent.MySQL {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		"",
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	// testDB := viper.GetString("mysql.db-name") + "_shadow"
	testDB := fmt.Sprintf("%s_shadow_%s", viper.GetString("mysql.db-name"), uuid.New().String()[0:8])

	assert.NoError(t, db.Exec("DROP DATABASE IF EXISTS "+testDB).Error)
	assert.NoError(t, db.Exec("CREATE DATABASE IF NOT EXISTS "+testDB).Error)

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		testDB,
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	assert.NoError(t, err)
	assert.NoError(t, db.AutoMigrate(persistent.StockModel{}))

	return persistent.NewMySQLWithDB(db)
}

func TestMySQLStockRepository_UpdateStock_Race(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)

	var (
		ctx          = context.Background()
		testItem     = "test-race-item"
		initialStock = 100
	)

	err := db.Create(ctx, &persistent.StockModel{
		ProductID: testItem,
		Quantity:  int64(initialStock),
	})
	assert.NoError(t, err)

	repo := NewMySQLStockRepository(db)
	var wg sync.WaitGroup
	goroutines := 10

	for range goroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := repo.UpdateStock(ctx, []*entity.ItemWithQuantity{
				{ID: testItem, Quantity: 1},
			}, func(ctx context.Context, existing, query []*entity.ItemWithQuantity) ([]*entity.ItemWithQuantity, error) {
				var newItems []*entity.ItemWithQuantity
				for _, e := range existing {
					for _, q := range query {
						if e.ID == q.ID {
							newItems = append(newItems, &entity.ItemWithQuantity{
								ID:       e.ID,
								Quantity: e.Quantity - q.Quantity,
							})
							break
						}
					}
				}
				return newItems, nil
			})
			assert.NoError(t, err)
		}()
	}
	wg.Wait()

	res, err := db.BatchGetStockByProductIDs(ctx, []string{testItem})

	assert.NoError(t, err)
	assert.NotEmpty(t, res, "res can not be empty")

	expected := initialStock - goroutines

	assert.Equal(t, int64(expected), res[0].Quantity)
}

func TestMySQLStockRepository_UpdateStock_OverSell(t *testing.T) {
	t.Parallel()
	db := setupTestDB(t)

	var (
		ctx          = context.Background()
		testItem     = "test-over-sell-item"
		initialStock = 10
	)

	err := db.Create(ctx, &persistent.StockModel{
		ProductID: testItem,
		Quantity:  int64(initialStock),
	})
	assert.NoError(t, err)

	repo := NewMySQLStockRepository(db)
	var wg sync.WaitGroup
	goroutines := 50

	for range goroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := repo.UpdateStock(ctx, []*entity.ItemWithQuantity{
				{ID: testItem, Quantity: 1},
			}, func(ctx context.Context, existing, query []*entity.ItemWithQuantity) ([]*entity.ItemWithQuantity, error) {
				var newItems []*entity.ItemWithQuantity
				for _, e := range existing {
					for _, q := range query {
						if e.ID == q.ID {
							newItems = append(newItems, &entity.ItemWithQuantity{
								ID:       e.ID,
								Quantity: e.Quantity - q.Quantity,
							})
							break
						}
					}
				}
				return newItems, nil
			})
			assert.NoError(t, err)
		}()
	}
	wg.Wait()

	res, err := db.BatchGetStockByProductIDs(ctx, []string{testItem})

	assert.NoError(t, err)
	assert.NotEmpty(t, res, "res can not be empty")

	assert.GreaterOrEqual(t, res[0].Quantity, int64(0))
}
