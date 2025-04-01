package integration

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/product"
)

type StripeAPI struct {
	apiKey string
}

func NewStripeAPI() *StripeAPI {
	logrus.Info("[NewStripeAPI] stripe-key found.")
	stripeKey := viper.GetString("stripe-key")
	if stripeKey == "" {
		logrus.Fatal("stripe-key is empty")
	}

	return &StripeAPI{
		apiKey: stripeKey,
	}
}

func (s *StripeAPI) GetPriceByProductID(ctx context.Context, pid string) (string, error) {
	// TODO: Logging

	stripe.Key = s.apiKey

	result, err := product.Get(pid, &stripe.ProductParams{})
	logrus.Debugf("PID: %s, result: %+v", pid, result)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}
