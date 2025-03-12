package discovery

import (
	"context"
	"github.com/FacundoChan/gorder-v1/common/discovery/consul"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

func RegisterToConsul(ctx context.Context, serviceName string) (func() error, error) {
	registry, err := consul.New(viper.GetString("consul.addr"))
	if err != nil {
		return func() error { return nil }, err
	}
	instanceID := GenerateInstanceID(serviceName)
	hostPort := viper.Sub(serviceName).GetString("grpc-addr")
	if err := registry.Register(ctx, instanceID, serviceName, hostPort); err != nil {
		return func() error { return nil }, err
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				logrus.Panicf("no heartbeat from %s to registery, err: %v", instanceID, err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	logrus.WithFields(logrus.Fields{
		"serviceName": serviceName,
		"addr":        hostPort,
	}).Info("Registered to Consul")

	return func() error {
		return registry.Deregister(ctx, instanceID, serviceName)
	}, nil
}
