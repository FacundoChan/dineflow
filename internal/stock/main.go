package main

import (
	"context"
	"github.com/FacundoChan/gorder-v1/common/config"
	"github.com/FacundoChan/gorder-v1/common/discovery"
	"github.com/FacundoChan/gorder-v1/common/genproto/stockpb"
	"github.com/FacundoChan/gorder-v1/common/logging"
	"github.com/FacundoChan/gorder-v1/common/server"
	"github.com/FacundoChan/gorder-v1/stock/ports"
	"github.com/FacundoChan/gorder-v1/stock/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	logging.Init()
	if err := config.NewViperConfig(); err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")
	logrus.Debugf("serviceName: %v, serverType: %v\n", serviceName, serverType)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := service.NewApplication(ctx)

	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = deregisterFunc()
	}()

	switch serverType {
	case "grpc":
		server.RunGRPCServer(serviceName, func(s *grpc.Server) {
			stockpb.RegisterStockServiceServer(s, ports.NewGRPCServer(app))
		})
	case "http":
		// TODO
		panic("http method not implement")
	default:
		panic("invalid stock server type")

	}

}
