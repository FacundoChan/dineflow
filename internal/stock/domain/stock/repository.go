package stock

import (
	"context"
	"fmt"
	"strings"
)

type Repository interface {
	GetItems(ctx context.Context, ids []string)
}

type NotFoundError struct {
	MissingIDs []string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("stock repository: not found: %s", strings.Join(e.MissingIDs, ","))
}
