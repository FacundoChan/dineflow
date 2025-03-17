package stock

import (
	"context"
	"fmt"
	"strings"

	"github.com/FacundoChan/gorder-v1/common/genproto/orderpb"
)

type Repository interface {
	GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error)
}

type NotFoundError struct {
	MissingIDs []string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("stock repository: not found: %s", strings.Join(e.MissingIDs, ","))
}
