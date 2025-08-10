package integration

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/product"
)

type StripeAPI struct {
	apiKey string
}

func NewStripeAPI(stripeKey string) *StripeAPI {
	return &StripeAPI{
		apiKey: stripeKey,
	}
}

// Deprecated: remote Stripe lookup is expensive. Prefer reading local price from DB via repository.
func (s *StripeAPI) GetPriceByProductID(ctx context.Context, pid string) (string, error) {
	logrus.Warn("GetPriceByProductID via Stripe is deprecated; use local price from DB")
	stripe.Key = s.apiKey
	result, err := product.Get(pid, &stripe.ProductParams{})
	if err != nil {
		return "", err
	}
	return result.ID, nil
}
