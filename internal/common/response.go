package common

import (
	"net/http"

	"github.com/FacundoChan/gorder-v1/common/tracing"
	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
}

type response struct {
	ErrorNo int    `json:"errorno"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	TraceID string `json:"trace_id"`
}

func (base *BaseResponse) Response(ctx *gin.Context, err error, data any) {
	if err != nil {
		base.error(ctx, err)
	} else {
		base.success(ctx, data)
	}
}

func (base *BaseResponse) success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, response{
		ErrorNo: 0,
		Message: "success",
		Data:    data,
		TraceID: tracing.TraceID(ctx.Request.Context()),
	})
}

func (base *BaseResponse) error(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, response{
		ErrorNo: 2,
		Message: err.Error(),
		Data:    nil,
		TraceID: tracing.TraceID(ctx.Request.Context()),
	})
}
