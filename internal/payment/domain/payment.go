package domain

import (
	"context"

	"github.com/FacundoChan/gorder-v1/common/genproto/orderpb"
)

type Processor interface {
	CreatePaymentLink(context.Context, *orderpb.Order) (string, error)
}
