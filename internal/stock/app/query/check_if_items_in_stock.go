package query

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/FacundoChan/dineflow/common/decorator"
	"github.com/FacundoChan/dineflow/common/handler/redis"
	"github.com/FacundoChan/dineflow/common/logging"
	domainStripe "github.com/FacundoChan/dineflow/stock/domain"
	domain "github.com/FacundoChan/dineflow/stock/domain/stock"
	"github.com/FacundoChan/dineflow/stock/entity"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	redisLockPrefix = "check_stock_"
)

type CheckIfItemsInStock struct {
	Items []*entity.ItemWithQuantity
}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*entity.Item]

type checkIfItemsInStockHandler struct {
	stockRepo domain.Repository
	stripeAPI domainStripe.StripeService
}

func NewCheckIfItemsInStockHandler(
	stockRepo domain.Repository,
	stripeAPI domainStripe.StripeService,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) CheckIfItemsInStockHandler {
	if stockRepo == nil {
		panic("nil orderRepo")
	}
	if stripeAPI == nil {
		panic("nil stripeAPI")
	}
	return decorator.ApplyQueryDecorators[CheckIfItemsInStock, []*entity.Item](
		checkIfItemsInStockHandler{
			stockRepo: stockRepo,
			stripeAPI: stripeAPI,
		},
		logger,
		metricsClient,
	)
}

func (c checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*entity.Item, error) {
	logrus.Debugf("redis_lock_key: %s", getLockKey(query))
	if err := lock(ctx, getLockKey(query)); err != nil {
		return nil, errors.Wrapf(err, "redis lock error: key=%s", getLockKey(query))
	}

	defer func() {
		if err := unlock(ctx, getLockKey(query)); err != nil {
			logging.Warnf(ctx, nil, "redis unlock failed, err=%v", err)
		}
	}()

	// Use local DB price via repository instead of calling Stripe per item
	// Read items' name and price from stock repository
	var res []*entity.Item
	ids := make([]string, 0, len(query.Items))
	for _, it := range query.Items {
		ids = append(ids, it.ID)
	}
	pricedItems, err := c.stockRepo.GetItems(ctx, ids)
	if err != nil {
		return nil, err
	}
	// Map quantities back
	qty := make(map[string]int32)
	for _, q := range query.Items {
		qty[q.ID] += q.Quantity
	}
	for _, pi := range pricedItems {
		res = append(res, &entity.Item{
			ID:       pi.ID,
			Name:     pi.Name,
			Quantity: qty[pi.ID],
			Price:    pi.Price,
		})
	}

	// Update the stock
	if err := c.checkStock(ctx, query.Items); err != nil {
		return nil, err
	}

	return res, nil
}

func getLockKey(query CheckIfItemsInStock) string {
	var ids []string
	for _, i := range query.Items {
		ids = append(ids, i.ID)
	}
	// Sort the ids to avoid different locks for the same items
	sort.Strings(ids)
	return redisLockPrefix + strings.Join(ids, "_")
}

func lock(ctx context.Context, key string) error {
	_, err := redis.SetNX(ctx, redis.LocalClient(), key, "1", 5*time.Minute)
	return err
}

func unlock(ctx context.Context, key string) error {
	return redis.Del(ctx, redis.LocalClient(), key)
}

func (c checkIfItemsInStockHandler) checkStock(ctx context.Context, queryItems []*entity.ItemWithQuantity) error {
	logrus.Debug("checkStock called")
	var ids []string
	for _, item := range queryItems {
		ids = append(ids, item.ID)
	}
	records, err := c.stockRepo.GetStock(ctx, ids)
	if err != nil {
		return err
	}
	idQuantityMap := make(map[string]int32)
	for _, record := range records {
		idQuantityMap[record.ID] += record.Quantity
	}
	var (
		ok             = true
		notEnoughItems []struct {
			ID   string
			Want int32
			Have int32
		}
	)
	for _, item := range queryItems {
		if item.Quantity > idQuantityMap[item.ID] {
			ok = false
			notEnoughItems = append(notEnoughItems, struct {
				ID   string
				Want int32
				Have int32
			}{
				ID:   item.ID,
				Want: item.Quantity,
				Have: idQuantityMap[item.ID],
			})
		}
	}

	if ok {
		c.stockRepo.UpdateStock(ctx, queryItems, func(ctx context.Context,
			existing []*entity.ItemWithQuantity,
			query []*entity.ItemWithQuantity,
		) ([]*entity.ItemWithQuantity, error) {
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
		return nil
	}

	return domain.ExceedStockError{NotEnoughItems: notEnoughItems}
}
