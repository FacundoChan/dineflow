package service

import (
	"context"
	"github.com/FacundoChan/gorder-v1/order/adapters"
	"github.com/FacundoChan/gorder-v1/order/app"
)

func NewApplication(ctx context.Context) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	return app.Application{
		Queries: New(orderRepo),
	}
}
