package redis

import (
	"context"
	"errors"
	"time"

	"github.com/FacundoChan/dineflow/common/logging"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func SetNX(ctx context.Context, client *redis.Client, key, value string, ttl time.Duration) (ok bool, err error) {
	now := time.Now()
	defer func() {
		l := logrus.WithContext(ctx).WithFields(logrus.Fields{
			"start":       now,
			"key":         key,
			"value":       value,
			logging.Error: err,
			logging.Cost:  time.Since(now).Milliseconds(),
		})
		if err == nil {
			l.Info("_redis_setnx_success")
		} else {
			l.Info("_redis_setnx_error")
		}
	}()

	if client == nil {
		return false, errors.New("redis client is nil")
	}

	ok, err = client.SetNX(ctx, key, value, ttl).Result()

	return ok, err
}

func Del(ctx context.Context, client *redis.Client, key string) (err error) {
	now := time.Now()
	defer func() {
		l := logrus.WithContext(ctx).WithFields(logrus.Fields{
			"start":       now,
			"key":         key,
			logging.Error: err,
			logging.Cost:  time.Since(now).Milliseconds(),
		})
		if err == nil {
			l.Info("_redis_del_success")
		} else {
			l.Info("_redis_del_error")
		}
	}()

	if client == nil {
		return errors.New("redis client is nil")
	}

	_, err = client.Del(ctx, key).Result()

	return err
}
