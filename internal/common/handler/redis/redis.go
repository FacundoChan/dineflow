package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/FacundoChan/dineflow/common/handler/factory"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	configName    = "redis"
	localSupplier = "local"
)

var (
	singleton = factory.NewSingleton(supplier)
)

func Init() {
	conf := viper.GetStringMap(configName)
	for supplyName := range conf {
		Client(supplyName)
	}
}

func LocalClient() *redis.Client {
	return Client(localSupplier)
}

func Client(name string) *redis.Client {
	return singleton.Get(name).(*redis.Client)
}

func supplier(key string) any {
	confKey := configName + "." + key
	type Section struct {
		IP           string        `mapstructure:"ip"`
		Port         string        `mapstructure:"port"`
		PoolSize     int           `mapstructure:"pool_size"`
		MaxConn      int           `mapstructure:"max_conn"`
		ConnTimeout  time.Duration `mapstructure:"conn_timeout"`
		ReadTimeout  time.Duration `mapstructure:"read_timeout"`
		WriteTimeout time.Duration `mapstructure:"write_timeout"`
	}

	var c Section
	if err := viper.UnmarshalKey(confKey, &c); err != nil {
		logrus.WithFields(logrus.Fields{
			"confKey": confKey,
			"err":     err,
		}).Error("Failed to unmarshal redis config")
		panic(err)
	}
	client := redis.NewClient(&redis.Options{
		Network:         "tcp",
		Addr:            fmt.Sprintf("%s:%s", c.IP, c.Port),
		ReadTimeout:     c.ReadTimeout * time.Millisecond,
		WriteTimeout:    c.WriteTimeout * time.Millisecond,
		// PoolSize:        c.PoolSize,
		// MaxActiveConns:  c.MaxConn,
		ConnMaxLifetime: c.ConnTimeout * time.Millisecond,
	})
	// Ping the Redis server to check if it's available
	if err := client.Ping(context.Background()).Err(); err != nil {
		logrus.WithFields(logrus.Fields{
			"addr": client.Options().Addr,
			"err":  err,
		}).Error("Redis server is not available")
	} else {
		logrus.WithFields(logrus.Fields{
			"addr": client.Options().Addr,
			// "version": client.Info(context.Background(), "server").String(),
		}).Info("Redis server is available")
	}
	return client
}
