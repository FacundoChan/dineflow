package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/FacundoChan/dineflow/common/broker"
	"github.com/FacundoChan/dineflow/common/genproto/orderpb"
	"github.com/FacundoChan/dineflow/common/logging"
	"github.com/FacundoChan/dineflow/payment/domain"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
	"go.opentelemetry.io/otel"
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
	logrus.WithContext(ctx.Request.Context()).Info("Got webhook from stripe")
	var err error

	defer func() {
		if err != nil {
			logging.Warnf(ctx.Request.Context(), nil, "handleWebHook err=%v", err)
		} else {
			logging.Infof(ctx.Request.Context(), nil, "handleWebHook success")
		}
	}()

	const MaxBodyBytes = int64(65536)
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		err = errors.Wrap(err, "Error reading request body")
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
		err = errors.Wrap(err, "Error verifying webhook signature")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			err = errors.Wrap(err, "Error unmarshalling checkout session")
			ctx.JSON(http.StatusBadRequest, err)
			return
		}

		if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
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
				err = errors.Wrap(err, "Error marshalling order")
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// TODO: MQ logging
			tr := otel.Tracer("rabbit-mq")
			mqCtx, span := tr.Start(ctx.Request.Context(), fmt.Sprintf("rabbit-mq.%s.publish", broker.EventOrderPaid))
			defer span.End()

			headers := broker.InjectRabbitMQHeaders(mqCtx)

			err = h.channel.PublishWithContext(mqCtx, broker.EventOrderPaid, "", false, false, amqp.Publishing{
				ContentType:  "application/json",
				DeliveryMode: amqp.Persistent,
				Body:         marshalledOrder,
				Headers:      headers,
			})
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			logrus.WithContext(ctx).Infof("message published to %s, body: %s", broker.EventOrderPaid, string(marshalledOrder))
			ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	default:
		logrus.WithContext(ctx).Infof("Unhandled event type: %s\n", event.Type)
	}

	// Ensure we end the handler only once; if nothing wrote, return 200
	if !ctx.IsAborted() && ctx.Writer.Status() == 200 && ctx.Writer.Size() <= 0 {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
