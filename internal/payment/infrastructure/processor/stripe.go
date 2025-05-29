package processor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/FacundoChan/dineflow/common/genproto/orderpb"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/product"
)

type StripeProcessor struct {
	apiKey string
}

func NewStripeProcessor(apiKey string) *StripeProcessor {
	if apiKey == "" {
		panic("empty api key")
	}
	stripe.Key = apiKey
	return &StripeProcessor{apiKey: apiKey}
}

const (
	// TODO: update successURL
	successURL = "http://localhost:3000/success"
)

func (s StripeProcessor) CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error) {
	logrus.Debugf("CreatePaymentLink %+v", order)
	var items []*stripe.CheckoutSessionLineItemParams

	for _, item := range order.Items {
		logrus.Debugf("adding item %+v", item)
		priceID, err := s.GetPriceByProductID(ctx, item.ID)
		if err != nil {
			logrus.Errorf("ID: %s not found in stripe", item.ID)
		}
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(priceID),
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}
	marshelledItems, _ := json.Marshal(order.Items)
	metadata := map[string]string{
		"orderID":     order.ID,
		"customerID":  order.CustomerID,
		"status":      order.Status,
		"paymentLink": order.PaymentLink,
		"items":       string(marshelledItems),
	}
	logrus.Debugf("metadata generated from StripeProccessor %+v", metadata)

	params := &stripe.CheckoutSessionParams{
		Metadata:   metadata,
		LineItems:  items,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(fmt.Sprintf("%s?customerID=%s&orderID=%s", successURL, order.CustomerID, order.ID)),
	}
	result, err := session.New(params)
	// logrus.Debugf("result generated from StripeProccessor %+v", result)
	if err != nil {
		return "", err
	}
	return result.URL, nil
}

func (s StripeProcessor) GetPriceByProductID(ctx context.Context, pid string) (string, error) {
	// TODO: Logging

	stripe.Key = s.apiKey

	result, err := product.Get(pid, &stripe.ProductParams{})
	logrus.Debugf("PID: %s, result: %+v\n", pid, result)
	logrus.Debugf("PID: %s, result.DefaultPrice.ID: %+v\n", pid, result.DefaultPrice.ID)
	if err != nil {
		return "", err
	}
	return result.DefaultPrice.ID, nil
}
