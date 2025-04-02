package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/FacundoChan/gorder-v1/common/broker"
	"github.com/FacundoChan/gorder-v1/common/client"
	_ "github.com/FacundoChan/gorder-v1/common/config"
	"github.com/FacundoChan/gorder-v1/common/logging"
	"github.com/FacundoChan/gorder-v1/common/tracing"
	"github.com/FacundoChan/gorder-v1/kitchen/adapters"
	"github.com/FacundoChan/gorder-v1/kitchen/infrastructure/consumer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	logging.Init()
}

func main() {
	serviceName := viper.GetString("kitchen.service-name")
	serverType := viper.GetString("kitchen.server-to-run")
	logrus.Debugf("serviceName: %v, serverType: %v\n", serviceName, serverType)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)

	client, closeClient, err := client.NewOrderGRPCClient(ctx)
	if err != nil {
		logrus.Fatal(err)
	}
	defer closeClient()
	orderGPRC := adapters.NewOrderGRPC(client)

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

	go consumer.NewConsumer(orderGPRC).Listen(ch)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigs
		logrus.Infof("Receive exit signal, exiting ...")
		os.Exit(0)
	}()

	logrus.Printf("To exit, press Ctrl + C")

	select {}

}
