package processor

import (
	"context"
	"encoding/json"
	"fmt"

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

const (
	// TODO: update successURL
	successURL = "http://localhost:5173/success"
)

func (s StripeProcessor) CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error) {
	logrus.Debugf("CreatePaymentLink %+v", order)
	var items []*stripe.CheckoutSessionLineItemParams

	for _, item := range order.Items {
		// logrus.Debugf("adding item %+v", item)
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			// TODO: Price
			//Price:    stripe.String(item.PriceID),
			Price:    stripe.String("price_1R1hXmDqhAs8dvRuCqn7mlgS"),
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
