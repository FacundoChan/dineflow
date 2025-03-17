package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/FacundoChan/gorder-v1/common/broker"
	"github.com/FacundoChan/gorder-v1/common/genproto/orderpb"
	"github.com/FacundoChan/gorder-v1/payment/domain"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
	"io"
	"net/http"
	"os"
)

type PaymentHandler struct {
	channel *amqp.Channel
}

func NewPaymentHandler(ch *amqp.Channel) *PaymentHandler {
	return &PaymentHandler{channel: ch}
}

func (h *PaymentHandler) RegisterRoutes(c *gin.Engine) {
	c.POST("/api/webhook", h.handleWebHook)

}

func (h *PaymentHandler) handleWebHook(ctx *gin.Context) {
	logrus.Info("Got webhook from stripe")
	const MaxBodyBytes = int64(65536)
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		ctx.Writer.WriteHeader(http.StatusServiceUnavailable)
		ctx.JSON(http.StatusServiceUnavailable, err)
		return
	}

	// Pass the request body and Stripe-Signature header to ConstructEvent, along
	// with the webhook signing key.
	//event, err := webhook.ConstructEvent(payload, ctx.Request.Header.Get("Stripe-Signature"),
	//	viper.GetString("endpoint-stripe-secret"))
	event, err := webhook.ConstructEventWithOptions(payload, ctx.Request.Header.Get("Stripe-Signature"),
		viper.GetString("endpoint-stripe-secret"), webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		})

	if err != nil {
		logrus.Infof("Error verifying webhook signature: %v\n", err)
		ctx.Writer.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			logrus.Infof("Error unmarshalling checkout session: %v\n", err)
			ctx.JSON(http.StatusBadRequest, err)
			return
		}

		if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
			logrus.Infof("payment for checkout session %s is paid", session.ID)
			logrus.Debugf("session metadata: %+v", session.Metadata)

			_, cancel := context.WithCancel(context.TODO())
			defer cancel()

			var items []*orderpb.Item
			_ = json.Unmarshal([]byte(session.Metadata["items"]), &items)
			marshalledOrder, err := json.Marshal(&domain.Order{
				ID:          session.Metadata["orderID"],
				CustomerID:  session.Metadata["customerID"],
				Status:      string(stripe.CheckoutSessionPaymentStatusPaid),
				PaymentLink: session.Metadata["paymentLink"],
				Items:       items,
			})
			if err != nil {
				logrus.Infof("Error marshalling order: %v\n", err)
				ctx.JSON(http.StatusBadRequest, err.Error())
				return
			}
			err = h.channel.PublishWithContext(ctx, broker.EventOrderPaid, "", false, false, amqp.Publishing{
				ContentType:  "application/json",
				DeliveryMode: amqp.Persistent,
				Body:         marshalledOrder,
			})
			if err != nil {
				ctx.JSON(http.StatusBadRequest, err.Error())
				return
			}
			logrus.Infof("message published to %s, body: %s", broker.EventOrderPaid, string(marshalledOrder))
			ctx.JSON(http.StatusOK, nil)
		}
	default:
		logrus.Infof("Unhandled event type: %s\n", event.Type)
	}

	ctx.Writer.WriteHeader(http.StatusOK)
}
