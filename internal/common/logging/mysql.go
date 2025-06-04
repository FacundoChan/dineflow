package logging

import (
	"context"
	"strings"
	"time"

	format "github.com/FacundoChan/dineflow/common/format"
	"github.com/sirupsen/logrus"
)

const (
	Method   = "method"
	Args     = "args"
	Cost     = "cost_ms"
	Response = "response"
	Error    = "err"
)

type ArgFormatter interface {
	FormatArg() (string, error)
}

func WhenMySQL(ctx context.Context, method string, args ...any) (logrus.Fields, func(any, *error)) {
	fields := logrus.Fields{
		Method: method,
		Args:   formatArgs(args),
	}

	start := time.Now()
	return fields, func(response any, err *error) {
		level, msg := logrus.InfoLevel, "mysql_success"
		fields[Cost] = time.Since(start).Microseconds()
		fields[Response] = response

		if err != nil && (*err != nil) {
			level, msg = logrus.ErrorLevel, "mysql_error"
			fields[Error] = (*err).Error()
		}

		logf(ctx, level, fields, "%s", msg)
	}

}

func formatArgs(args []any) string {
	var item []string
	for _, arg := range args {
		item = append(item, formatArg(arg))

	}
	return strings.Join(item, ", ")
}

func formatArg(arg any) (str string) {
	var err error

	defer func() {
		if err != nil {
			str = "unsupported type in formatMySQLArg, err=" + err.Error()
		}
	}()

	switch v := arg.(type) {
	case ArgFormatter:
		str, err = v.FormatArg()
	default:
		str, err = format.MarshalString(v)
	}

	return str
}
