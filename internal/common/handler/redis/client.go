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

func SetEX(ctx context.Context, client *redis.Client, key, value string, expiration time.Duration) (err error) {
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
			l.Info("_redis_setex_success")
		} else {
			l.Info("_redis_setex_error")
		}
	}()

	if client == nil {
		return errors.New("redis client is nil")
	}

	_, err = client.SetEx(ctx, key, value, expiration).Result()

	return err
}

func GetEX(ctx context.Context, client *redis.Client, key string, expiration time.Duration) (value string, err error) {
	now := time.Now()
	defer func() {
		l := logrus.WithContext(ctx).WithFields(logrus.Fields{
			"start":       now,
			"key":         key,
			logging.Error: err,
			logging.Cost:  time.Since(now).Milliseconds(),
		})
		if err == nil {
			l.Info("_redis_getex_success")
		} else if err == redis.Nil {
			l.Info("_redis_getex_nil")
		} else {
			l.Info("_redis_getex_error")
		}
	}()

	if client == nil {
		return "", errors.New("redis client is nil")
	}

	value, err = client.GetEx(ctx, key, expiration).Result()

	return value, err
}
