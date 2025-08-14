package broker

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"

	// Import config to trigger init() for config loading
	_ "github.com/FacundoChan/dineflow/common/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	DLX                = "dlx"
	DLQ                = "dlq"
	amqpRetryHeaderKey = "x-retry-count"
)

var (
	maxRetryCount = viper.GetInt64("rabbit-mq.max-retry")
)

func Connect(user, password, host, port string) (*amqp.Channel, func() error) {
	conn, err := Dial(user, password, host, port)
	if err != nil {
		logrus.Fatal(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatal(err)
	}
	if err := InitChannel(ch); err != nil {
		logrus.Fatal(err)
	}
	return ch, conn.Close
}

// Dial establishes and returns a RabbitMQ connection without opening a channel.
func Dial(user, password, host, port string) (*amqp.Connection, error) {
	addr := fmt.Sprintf("amqp://%s:%s@%s:%s", user, password, host, port)
	return amqp.Dial(addr)
}

// InitChannel declares required exchanges/queues on the given channel.
func InitChannel(ch *amqp.Channel) error {
	if err := ch.ExchangeDeclare(EventOrderCreated, "direct", true, false, false, false, nil); err != nil {
		return err
	}
	if err := ch.ExchangeDeclare(EventOrderPaid, "direct", true, false, false, false, nil); err != nil {
		return err
	}
	if err := createDLX(ch); err != nil {
		return err
	}
	return nil
}

func createDLX(ch *amqp.Channel) error {
	q, err := ch.QueueDeclare("shared_queue", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = ch.ExchangeDeclare(DLX, "fanout", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = ch.QueueBind(q.Name, "", DLX, false, nil)
	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(DLQ, true, false, false, false, nil)
	return err
}

func HandleRetry(ctx context.Context, ch *amqp.Channel, d *amqp.Delivery) error {
	if d.Headers == nil {
		d.Headers = amqp.Table{}
	}

	logrus.Debug("maxRetryCount: ", maxRetryCount)
	retryCount, ok := d.Headers[amqpRetryHeaderKey].(int64)
	if !ok {
		retryCount = 0
	}
	retryCount++
	d.Headers[amqpRetryHeaderKey] = retryCount
	if retryCount > maxRetryCount {
		logrus.Infof("moving msg %s to dlq", d.MessageId)
		return ch.PublishWithContext(ctx, "", DLQ, false, false, amqp.Publishing{
			Headers:      d.Headers,
			ContentType:  "application/json",
			Body:         d.Body,
			DeliveryMode: amqp.Persistent,
		})
	}
	logrus.Infof("retring msg %s, cnt=%d", d.MessageId, retryCount)
	time.Sleep(time.Second * time.Duration(retryCount))

	return ch.PublishWithContext(ctx, d.Exchange, d.RoutingKey, false, false, amqp.Publishing{
		Headers:      d.Headers,
		ContentType:  "application/json",
		Body:         d.Body,
		DeliveryMode: amqp.Persistent,
	})
}

type RabbitMQHeaderCarrier map[string]interface{}

func (r RabbitMQHeaderCarrier) Get(key string) string {
	value, ok := r[key]
	if !ok {
		return ""
	}
	return value.(string)
}

func (r RabbitMQHeaderCarrier) Set(key string, value string) {
	r[key] = value
}

func (r RabbitMQHeaderCarrier) Keys() []string {
	keys := make([]string, len(r))
	i := 0
	for key := range r {
		keys[i] = key
		i++
	}
	return keys
}

func InjectRabbitMQHeaders(ctx context.Context) map[string]interface{} {
	carrier := make(RabbitMQHeaderCarrier)
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	return carrier
}

func ExtractRabbitMQHeaders(ctx context.Context, headers map[string]interface{}) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, RabbitMQHeaderCarrier(headers))
}
