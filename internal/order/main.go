package main

import (
	"github.com/FacundoChan/gorder-v1/common/genproto/orderpb"
	"google.golang.org/grpc"
	"log"

	"github.com/FacundoChan/gorder-v1/common/config"
	"github.com/FacundoChan/gorder-v1/common/server"
	"github.com/FacundoChan/gorder-v1/order/ports"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	serviceName := viper.GetString("order.service-name")

	go server.RunGRPCServer(serviceName, func(s *grpc.Server) {
		svc := ports.NewGRPCServer()
		orderpb.RegisterOrderServiceServer(s, svc)
	})

	server.RunHTTPServer(serviceName, func(r *gin.Engine) {
		ports.RegisterHandlersWithOptions(r, HTTPServer{}, ports.GinServerOptions{
			BaseURL:      "api",
			Middlewares:  nil,
			ErrorHandler: nil,
		})
	})
}
