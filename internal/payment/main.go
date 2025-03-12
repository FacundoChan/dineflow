package main

import (
	"github.com/FacundoChan/gorder-v1/common/config"
	"github.com/FacundoChan/gorder-v1/common/logging"
	"github.com/FacundoChan/gorder-v1/common/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	logging.Init()
	if err := config.NewViperConfig(); err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	serviceName := viper.GetString("payment.service-name")
	serverType := viper.GetString("payment.server-to-run")
	paymentHandler := NewPaymentHandler()
	switch serverType {
	case "http":
		server.RunHTTPServer(serviceName, paymentHandler.RegisterRoutes)
	case "grpc":
		logrus.Panic("unsupported service type")
	default:
		logrus.Panic("unreachable code")
	}
}
