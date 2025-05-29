package main

import (
	"context"

	_ "github.com/FacundoChan/dineflow/common/config"
	"github.com/FacundoChan/dineflow/common/discovery"
	"github.com/FacundoChan/dineflow/common/genproto/stockpb"
	"github.com/FacundoChan/dineflow/common/logging"
	"github.com/FacundoChan/dineflow/common/server"
	"github.com/FacundoChan/dineflow/common/tracing"
	"github.com/FacundoChan/dineflow/stock/ports"
	"github.com/FacundoChan/dineflow/stock/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	logging.Init()
}

func main() {
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")
	logrus.Debugf("serviceName: %v, serverType: %v\n", serviceName, serverType)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)

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
