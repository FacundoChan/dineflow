package adapters

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/FacundoChan/gorder-v1/common/genproto/orderpb"
	domain "github.com/FacundoChan/gorder-v1/stock/domain/stock"
)

type MemoryStockRepository struct {
	lock  *sync.RWMutex
	store map[string]*orderpb.Item
}

var stub = map[string]*orderpb.Item{
	"item_id": {
		ID:       "foo_item",
		Name:     "bar_name",
		Quantity: 1000,
		PriceID:  "stub_item_price_id",
	},
	"item-id-1": {
		ID:       "item-id-1",
		Name:     "item-name-1",
		Quantity: 10,
		PriceID:  "item-price-1",
	},
	"item-id-2": {
		ID:       "item-id-2",
		Name:     "item-name-2",
		Quantity: 100,
		PriceID:  "item-price-2",
	},
	"item-id-3": {
		ID:       "item-id-3",
		Name:     "item-name-3",
		Quantity: 1000,
		PriceID:  "item-price-3",
	},
}

func NewMemoryStockRepository() *MemoryStockRepository {
	return &MemoryStockRepository{
		lock: &sync.RWMutex{},
		//store: make([]*stockpb.Item, 0),
		store: stub,
	}
}

func (m *MemoryStockRepository) GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	logrus.Debug("GetItems called")
	var (
		res     []*orderpb.Item
		missing []string
	)

	for _, id := range ids {
		if item, exist := m.store[id]; exist {
			res = append(res, item)
		} else {
			missing = append(missing, id)
		}
	}

	if len(res) == len(ids) {
		return res, nil
	}

	return res, domain.NotFoundError{
		MissingIDs: missing,
	}

}
