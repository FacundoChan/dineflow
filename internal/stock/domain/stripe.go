package domain

import "context"

type StripeService interface {
	GetPriceByProductID(ctx context.Context, pid string) (string, error)
}
