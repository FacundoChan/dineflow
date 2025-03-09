package decorator

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type MetricsClient interface {
	Inc(key string, value int)
}

// C: Command
// R: Result
type queryMetricsDecorator[C, R any] struct {
	client MetricsClient
	base   QueryHandler[C, R]
}

func (q queryMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	start := time.Now()
	actionName := strings.ToLower(generateActionName(cmd))

	defer func() {
		end := time.Since(start)
		q.client.Inc(fmt.Sprintf("query.%s.duration", actionName), int(end.Seconds()))
		if err == nil {
			q.client.Inc(fmt.Sprintf("query.%s.success", actionName), 1)
		} else {
			q.client.Inc(fmt.Sprintf("query.%s.fail", actionName), 1)
		}
	}()

	return q.base.Handle(ctx, cmd)
}
