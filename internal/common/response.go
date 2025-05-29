package common

import (
	"encoding/json"
	"net/http"

	"github.com/FacundoChan/dineflow/common/handler/errors"
	"github.com/FacundoChan/dineflow/common/tracing"
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
	errno, errmsg := errors.Output(nil)
	r := response{
		ErrorNo: errno,
		Message: errmsg,
		Data:    data,
		TraceID: tracing.TraceID(ctx.Request.Context()),
	}

	response, _ := json.Marshal(r)
	ctx.Set("response", response)
	ctx.JSON(http.StatusOK, r)
}

func (base *BaseResponse) error(ctx *gin.Context, err error) {
	errno, errmsg := errors.Output(err)
	r := response{
		ErrorNo: errno,
		Message: errmsg,
		Data:    nil,
		TraceID: tracing.TraceID(ctx.Request.Context()),
	}
	response, _ := json.Marshal(r)
	ctx.Set("response", response)
	ctx.JSON(http.StatusOK, r)
}
