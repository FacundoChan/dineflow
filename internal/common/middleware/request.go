package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/FacundoChan/dineflow/common/logging"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RequestLog(l *logrus.Entry) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestIn(ctx, l)
		defer requestOut(ctx, l)
		ctx.Next()
	}
}

func requestIn(ctx *gin.Context, l *logrus.Entry) {
	ctx.Set("request_start", time.Now())
	body := ctx.Request.Body
	bodyBytes, _ := io.ReadAll(body)
	ctx.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	var compactJSON bytes.Buffer
	_ = json.Compact(&compactJSON, bodyBytes)
	l.WithContext(ctx.Request.Context()).WithFields(logrus.Fields{
		"start":      time.Now().Unix(),
		logging.Args: compactJSON.String(),
		"from":       ctx.RemoteIP(),
		"uri":        ctx.Request.RequestURI,
	}).Info("_request_in")
}

func requestOut(ctx *gin.Context, l *logrus.Entry) {
	response, _ := ctx.Get("response")
	start, _ := ctx.Get("request_start")
	startTime := start.(time.Time)
	l.WithContext(ctx.Request.Context()).WithFields(logrus.Fields{
		logging.Cost:     time.Since(startTime).Milliseconds(),
		logging.Response: string(response.([]byte)),
	}).Info("_request_out")
}
