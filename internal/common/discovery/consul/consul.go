package consul

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
)

type Registry struct {
	client *api.Client
}

var (
	consulClient *Registry
	once         sync.Once
	initErr      error
)

func New(consulAddr string) (*Registry, error) {
	once.Do(func() {
		config := api.DefaultConfig()
		config.Address = consulAddr
		client, err := api.NewClient(config)
		if err != nil {
			initErr = err
			return
		}
		consulClient = &Registry{client: client}
	})

	if initErr != nil {
		return nil, initErr
	}
	return consulClient, nil
}

func (r *Registry) Register(ctx context.Context, instanceID, serviceName, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("invalid host port")
	}
	host := parts[0]
	port, _ := strconv.Atoi(parts[1])
	return r.client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      instanceID,
		Name:    serviceName,
		Port:    port,
		Address: host,
		Check: &api.AgentServiceCheck{
			CheckID:                        instanceID,
			TLSSkipVerify:                  false,
			TTL:                            "5s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "10s",
		},
	})

}

func (r *Registry) Deregister(ctx context.Context, instanceID, serviceName string) error {
	logrus.WithFields(logrus.Fields{
		"isntanceID":  instanceID,
		"serviceName": serviceName,
	}).Info("deregister from consul")
	return r.client.Agent().CheckDeregister(instanceID)
}

func (r *Registry) Discover(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}
	var ips []string
	for _, entry := range entries {
		ips = append(ips, entry.Service.Address+":"+strconv.Itoa(entry.Service.Port))
	}
	return ips, nil
}

func (r *Registry) HealthCheck(instanceID, serviceName string) error {
	return r.client.Agent().UpdateTTL(instanceID, "online", api.HealthPassing)
}
