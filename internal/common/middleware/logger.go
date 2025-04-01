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
			"request_uri":     ctx.Request.RequestURI,
			"remote_addr":     ctx.Request.RemoteAddr,
			"client_ip":       ctx.ClientIP(),
			"full_path":       ctx.FullPath(),
		}).Info("request_out")
	}
}
