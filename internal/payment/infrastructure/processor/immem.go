package processor

import (
	"context"

	"github.com/FacundoChan/gorder-v1/common/genproto/orderpb"
)

type InmemProcessor struct {
}

func NewInmemProcessor() *InmemProcessor {
	return &InmemProcessor{}
}

func (i InmemProcessor) CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error) {
	// TODO:
	return "inmem-payment-link", nil
}
