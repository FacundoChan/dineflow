package adapters

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/FacundoChan/dineflow/common/logging"
	domain "github.com/FacundoChan/dineflow/order/domain/order"
	"github.com/sirupsen/logrus"
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

func (m *MemoryOrderRepository) Create(ctx context.Context, order *domain.Order) (created *domain.Order, err error) {
	_, dLog := logging.WhenRequest(ctx, "MemoryOrderRepository.Create", map[string]any{
		"order": order,
	})
	defer dLog(created, &err)

	m.lock.Lock()
	defer m.lock.Unlock()
	created = &domain.Order{
		ID:          strconv.FormatInt(time.Now().Unix(), 10),
		CustomerID:  order.CustomerID,
		Status:      order.Status,
		PaymentLink: order.PaymentLink,
		Items:       order.Items,
	}
	m.store = append(m.store, created)
	if logrus.GetLevel() == logrus.DebugLevel {
		logrus.Debugln("order_immem_repo_create")
		for i, order := range m.store {
			logrus.Debugf("m.store[%d]: %+v\n", i, order)
			for i, o := range order.Items {
				logrus.Debugf("     order-item[%d]: %+v", i, o)
			}
		}
	}
	logrus.WithFields(logrus.Fields{
		"input_order":        order,
		"store_after_create": m.store,
	}).Info("memory_order_repo_create")
	return created, nil
}

func (m *MemoryOrderRepository) Get(ctx context.Context, id, customerID string) (order *domain.Order, err error) {
	_, dLog := logging.WhenRequest(ctx, "MemoryOrderRepository.Get", map[string]any{
		"id":          id,
		"customer_id": customerID,
	})
	defer dLog(order, &err)

	m.lock.RLock()
	defer m.lock.RUnlock()

	if logrus.GetLevel() == logrus.DebugLevel {
		logrus.Debugln("order_immem_repo_get")
		for i, order := range m.store {
			logrus.Debugf("m.store[%d]: %+v\n", i, order)
			for i, o := range order.Items {
				logrus.Debugf("     order-item[%d]: %+v", i, o)
			}
		}
	}

	for _, order := range m.store {
		if order.ID == id && order.CustomerID == customerID {
			logging.Infof(ctx, nil, "memory_order_repo_get found id=%s customerID=%s res=%v", id, order.CustomerID, *order)
			return order, nil
		}
	}

	logging.Infof(ctx, nil, "memory_order_repo_get NOT found id=%s customerID=%s", id, customerID)
	return nil, domain.NotFoundError{OrderID: id}
}

func (m *MemoryOrderRepository) Update(ctx context.Context, order *domain.Order, updateFunc func(context.Context, *domain.Order) (*domain.Order, error)) (err error) {
	_, dLog := logging.WhenRequest(ctx, "MemoryOrderRepository.Update", map[string]any{
		"order": order,
	})
	defer dLog(nil, &err)

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
