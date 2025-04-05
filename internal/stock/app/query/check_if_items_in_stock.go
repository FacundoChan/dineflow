package query

import (
	"context"

	"github.com/FacundoChan/gorder-v1/common/decorator"
	domain "github.com/FacundoChan/gorder-v1/stock/domain/stock"
	"github.com/FacundoChan/gorder-v1/stock/entity"
	"github.com/FacundoChan/gorder-v1/stock/infrastructure/integration"
	"github.com/sirupsen/logrus"
)

type CheckIfItemsInStock struct {
	Items []*entity.ItemWithQuantity
}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*entity.Item]

type checkIfItemsInStockHandler struct {
	stockRepo domain.Repository
	stripeAPI *integration.StripeAPI
}

func NewCheckIfItemsInStockHandler(
	stockRepo domain.Repository,
	stripeAPI *integration.StripeAPI,
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
	// TODO: CheckIfItemsInStock NOT DONE YET: Should be pulled from database or Stripe
	if err := c.checkStock(ctx, query.Items); err != nil {
		return nil, err
	}

	var res []*entity.Item
	for _, item := range query.Items {
		priceID, err := c.stripeAPI.GetPriceByProductID(ctx, item.ID)
		if err != nil || priceID == "" {
			logrus.Warnf("GetPriceByProductID error, item ID=%s, err=%v", item.ID, err)
			return nil, err
		}
		res = append(res, &entity.Item{
			ID:       item.ID,
			Quantity: item.Quantity,
			PriceID:  priceID,
		})
	}

	// TODO: update the stock
	return res, nil
}

func (c checkIfItemsInStockHandler) checkStock(ctx context.Context, queryItems []*entity.ItemWithQuantity) error {
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
		return nil
	}

	return domain.ExceedStockError{NotEnoughItems: notEnoughItems}
}
