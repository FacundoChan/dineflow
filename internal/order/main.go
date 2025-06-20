package main

import (
	"context"

	"github.com/FacundoChan/dineflow/common/broker"
	_ "github.com/FacundoChan/dineflow/common/config"
	"github.com/FacundoChan/dineflow/common/discovery"
	"github.com/FacundoChan/dineflow/common/genproto/orderpb"
	"github.com/FacundoChan/dineflow/common/logging"
	"github.com/FacundoChan/dineflow/common/server"
	"github.com/FacundoChan/dineflow/common/tracing"
	"github.com/FacundoChan/dineflow/order/infrastructure/consumer"
	"github.com/FacundoChan/dineflow/order/ports"
	"github.com/FacundoChan/dineflow/order/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	logging.Init()
}

func main() {
	serviceName := viper.GetString("order.service-name")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)

	app, cleanup := service.NewApplication(ctx)
	defer cleanup()

	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = deregisterFunc()
	}()

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

	go consumer.NewConsumer(app).Listen(ch)

	go server.RunGRPCServer(serviceName, func(s *grpc.Server) {
		svc := ports.NewGRPCServer(app)
		orderpb.RegisterOrderServiceServer(s, svc)
	})

	server.RunHTTPServer(serviceName, func(r *gin.Engine) {
		ports.RegisterHandlersWithOptions(r, HTTPServer{
			app: app,
		}, ports.GinServerOptions{
			BaseURL:      "api",
			Middlewares:  nil,
			ErrorHandler: nil,
		})
	})
}
