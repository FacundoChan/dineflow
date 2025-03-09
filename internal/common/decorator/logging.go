package decorator

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

// C: Command
// R: Result
type queryLoggingDecorator[C, R any] struct {
	logger *logrus.Entry
	base   QueryHandler[C, R]
}

func (q queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	logger := q.logger.WithFields(logrus.Fields{
		"query":      generateActionName(cmd),
		"query_body": fmt.Sprintf("%#v", cmd),
	})
	logger.Debug("Executing query")
	defer func() {
		if err == nil {
			logger.Info("Query succeeded")
		} else {
			logger.Error("Failed to execute query", err)
		}
	}()

	return q.base.Handle(ctx, cmd)
}

func generateActionName(cmd any) string {
	return strings.Split(fmt.Sprintf("%T", cmd), ".")[1]
}
