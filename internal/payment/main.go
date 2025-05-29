package main

import (
	"context"

	"github.com/FacundoChan/dineflow/common/broker"
	_ "github.com/FacundoChan/dineflow/common/config"
	"github.com/FacundoChan/dineflow/common/logging"
	"github.com/FacundoChan/dineflow/common/server"
	"github.com/FacundoChan/dineflow/common/tracing"
	"github.com/FacundoChan/dineflow/payment/infrastructure/consumer"
	"github.com/FacundoChan/dineflow/payment/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	logging.Init()
}

func main() {
	serviceName := viper.GetString("payment.service-name")
	serverType := viper.GetString("payment.server-to-run")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)

	application, cleanup := service.NewApplication(ctx)
	defer cleanup()

	ch, closeCh := broker.Connect(
		viper.GetString("rabbit-mq.user"),
		viper.GetString("rabbit-mq.password"),
		viper.GetString("rabbit-mq.host"),
		viper.GetString("rabbit-mq.port"),
	)
	logrus.Infof("Connected to rabbitmq")
	defer func() {
		_ = closeCh()
		_ = ch.Close()
	}()

	go consumer.NewConsumer(application).Listen(ch)

	paymentHandler := NewPaymentHandler(ch)

	switch serverType {
	case "http":
		server.RunHTTPServer(serviceName, paymentHandler.RegisterRoutes)
	case "grpc":
		logrus.Panic("unsupported service type")
	default:
		logrus.Panic("unreachable code")
	}
}
