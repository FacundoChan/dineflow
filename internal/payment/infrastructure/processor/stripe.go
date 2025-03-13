package processor

import (
	"context"
	"encoding/json"
	"github.com/FacundoChan/gorder-v1/common/genproto/orderpb"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
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

var (
	successURL = "http://localhost:8282/success"
)

func (s StripeProcessor) CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error) {
	var items []*stripe.CheckoutSessionLineItemParams

	for _, item := range order.Items {
		logrus.Debugf("adding item %+v", item)
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			//Price:    stripe.String(item.PriceID),
			Price:    stripe.String("price_1R1hXmDqhAs8dvRuCqn7mlgS"),
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}
	marshelledItems, _ := json.Marshal(order.Items)
	metadata := map[string]string{
		"orderID":    order.ID,
		"customerID": order.CustomerID,
		"status":     order.Status,
		"items":      string(marshelledItems),
	}

	params := &stripe.CheckoutSessionParams{
		Metadata:   metadata,
		LineItems:  items,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(successURL),
	}
	result, err := session.New(params)
	if err != nil {
		return "", err
	}
	return result.URL, nil
}
