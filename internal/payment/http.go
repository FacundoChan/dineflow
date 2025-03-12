package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PaymentHandler struct {
}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{}
}

func (h *PaymentHandler) RegisterRoutes(c *gin.Engine) {
	c.POST("/api/webhook", h.handleWebHook)

}

func (h *PaymentHandler) handleWebHook(ctx *gin.Context) {
	logrus.Info("Got webhook from stripe")
}
