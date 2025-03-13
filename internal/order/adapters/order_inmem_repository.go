package adapters

import (
	"context"
	domain "github.com/FacundoChan/gorder-v1/order/domain/order"
	"github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"time"
)

type MemoryOrderRepository struct {
	lock  *sync.RWMutex
	store []*domain.Order
}

func NewMemoryOrderRepository() *MemoryOrderRepository {
	s := make([]*domain.Order, 0)
	s = append(s, &domain.Order{
		ID:          "fake-ID",
		CustomerID:  "fake-customer-id",
		Status:      "fake-status",
		PaymentLink: "fake-payment-link",
		Items:       nil,
	})
	return &MemoryOrderRepository{
		lock: &sync.RWMutex{},
		//store: make([]*domain.Order, 0),
		store: s,
	}
}

func (m *MemoryOrderRepository) Create(_ context.Context, order *domain.Order) (*domain.Order, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	res := &domain.Order{
		ID:          strconv.FormatInt(time.Now().Unix(), 10),
		CustomerID:  order.CustomerID,
		Status:      order.Status,
		PaymentLink: order.PaymentLink,
		Items:       order.Items,
	}
	m.store = append(m.store, res)
	logrus.WithFields(logrus.Fields{
		"input_order":        order,
		"store_after_create": m.store,
	}).Info("memory_order_repo_create")
	return res, nil
}

func (m *MemoryOrderRepository) Get(_ context.Context, id, customerID string) (*domain.Order, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if logrus.GetLevel() == logrus.DebugLevel {
		for i, order := range m.store {
			logrus.Debugf("m.store[%d]: %+v", i, order)
		}
	}

	for _, order := range m.store {
		if order.ID == id && order.CustomerID == customerID {
			logrus.Debugf("memory_order_repo_get found id=%s customerID=%s res=%v", id, order.CustomerID, *order)
			return order, nil
		}
	}

	logrus.Debugf("memory_order_repo_get NOT found id=%s customerID=%s", id, customerID)
	return nil, domain.NotFoundError{OrderID: id}
}

func (m *MemoryOrderRepository) Update(ctx context.Context, order *domain.Order, updateFunc func(context.Context, *domain.Order) (*domain.Order, error)) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	found := false
	for i, o := range m.store {
		if o.ID == order.ID && o.CustomerID == order.CustomerID {
			found = true
			updatedOrder, err := updateFunc(ctx, order)
			if err != nil {
				return err
			}
			m.store[i] = updatedOrder
		}
	}

	if !found {
		return domain.NotFoundError{OrderID: order.ID}
	}

	return nil
}
