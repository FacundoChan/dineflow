package main

import (
	"github.com/FacundoChan/gorder-v1/common/genproto/stockpb"
	"github.com/FacundoChan/gorder-v1/common/server"
	"github.com/FacundoChan/gorder-v1/stock/ports"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")

	switch serverType {
	case "grpc":
		server.RunGRPCServer(serviceName, func(s *grpc.Server) {
			stockpb.RegisterStockServiceServer(s, ports.NewGRPCServer())
		})
	case "http":
		// TODO
		panic("http method not implement")
	default:
		panic("invalid stock server type")

	}

}
