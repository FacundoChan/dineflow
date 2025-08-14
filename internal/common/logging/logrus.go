package logging

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/FacundoChan/dineflow/common/tracing"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	// prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// Optional: Use logging.Info, Warnf ...
// Or: Add hook, Use logrus.Info, Warnf directly

func Init() {
	SetFormatter(logrus.StandardLogger())
	logrus.SetLevel(logrus.DebugLevel)
	setOutput(logrus.StandardLogger())
	logrus.AddHook(&traceHook{})
}

func setOutput(logger *logrus.Logger) {
	var (
		folder    = "./log/"
		filePath  = "dineflow.log"
		errorPath = "dineflow.error"
		linkName  = folder + "dineflow.log"
	)

	if err := os.MkdirAll(folder, 0750); err != nil && !os.IsExist(err) {
		panic(err)
	}
	file, err := os.OpenFile(folder+filePath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	logger.SetOutput(file)

	rotateInfo, err := rotatelogs.New(
		folder+filePath+".%Y%m%d"+".log",
		rotatelogs.WithLinkName(linkName),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)

	if err != nil {
		panic(err)
	}

	rotateError, err := rotatelogs.New(
		folder+errorPath+".%Y%m%d"+".log",
		rotatelogs.WithLinkName(folder+errorPath),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)

	if err != nil {
		panic(err)
	}
	rotationMap := lfshook.WriterMap{
		logrus.DebugLevel: rotateInfo,
		logrus.InfoLevel:  rotateInfo,
		logrus.WarnLevel:  rotateError,
		logrus.ErrorLevel: rotateError,
		logrus.FatalLevel: rotateError,
		logrus.PanicLevel: rotateError,
	}

	logrus.AddHook(lfshook.NewHook(rotationMap, &logrus.JSONFormatter{
		TimestampFormat: time.DateTime,
	}))
}

func SetFormatter(logger *logrus.Logger) {
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyMsg:   "message",
		},
	})
	if isLocal, _ := strconv.ParseBool(os.Getenv("LOCAL_ENV")); isLocal {
		// console output setting
		// logger.SetFormatter(&prefixed.TextFormatter{
		// 	TimestampFormat: time.RFC3339,
		// 	ForceColors:     true,
		// 	ForceFormatting: true,
		// })
	}

}

func InfofWithCost(ctx context.Context, fields logrus.Fields, start time.Time, format string, args ...any) {
	fields[Cost] = time.Since(start).Milliseconds()
	Infof(ctx, fields, format, args...)
}

func logf(ctx context.Context, level logrus.Level, fields logrus.Fields, format string, args ...any) {
	logrus.WithContext(ctx).WithFields(fields).Logf(level, format, args...)
}

func Infof(ctx context.Context, fields logrus.Fields, format string, args ...any) {
	logrus.WithContext(ctx).WithFields(fields).Infof(format, args...)
}

func Errorf(ctx context.Context, fields logrus.Fields, format string, args ...any) {
	logrus.WithContext(ctx).WithFields(fields).Errorf(format, args...)
}

func Warnf(ctx context.Context, fields logrus.Fields, format string, args ...any) {
	logrus.WithContext(ctx).WithFields(fields).Warnf(format, args...)
}

func Panicf(ctx context.Context, fields logrus.Fields, format string, args ...any) {
	logrus.WithContext(ctx).WithFields(fields).Panicf(format, args...)
}

type traceHook struct{}

// Fire implements logrus.Hook.
func (t *traceHook) Fire(entry *logrus.Entry) error {
	if entry.Context != nil {
		entry.Data["hook_trace"] = tracing.TraceID(entry.Context)
	}
	return nil
}

// Levels implements logrus.Hook.
func (t *traceHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
