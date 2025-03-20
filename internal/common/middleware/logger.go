package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StructuredLog(l *logrus.Entry) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		elapsed := time.Since(start)
		l.WithFields(logrus.Fields{
			"time_elapsed_ms": elapsed.Milliseconds(),
			"request_url":     ctx.Request.RequestURI,
			"client_ip":       ctx.ClientIP(),
		}).Info("request_out")
	}
}
