package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/AydarZaynutdinov/pow/internal/metrics"

	"github.com/AydarZaynutdinov/pow/internal/config"
	"github.com/redis/go-redis/v9"
)

const (
	defaultReconnectionTimeout = 5
)

type Cache struct {
	client  *redis.Client
	metrics *metrics.CacheMetrics
}

func NewCache(ctx context.Context, cfg config.Cache, metrics *metrics.CacheMetrics) (*Cache, error) {
	opts := &redis.Options{
		Addr:     cfg.Address,
		Username: cfg.Username,
		Password: cfg.Password,
		PoolSize: cfg.PoolSize,
	}

	client := redis.NewClient(opts)

	err := client.Ping(ctx).Err()
	if err == nil {
		return &Cache{
			client:  client,
			metrics: metrics,
		}, nil
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	timeoutExceeded := time.After(time.Second * time.Duration(defaultReconnectionTimeout))

	for {
		select {
		case <-timeoutExceeded:
			return nil, fmt.Errorf("connection failed after %d timeout", defaultReconnectionTimeout)
		case <-ticker.C:
			err = client.Ping(ctx).Err()
			if err == nil {
				return &Cache{
					client:  client,
					metrics: metrics,
				}, nil
			}
		}
	}
}

func (c *Cache) Exists(ctx context.Context, key string) (int64, error) {
	success, fail := c.observe("EXISTS")
	defer fail()

	res, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		fail()
	} else {
		success()
	}

	return res, err
}

func (c *Cache) Set(ctx context.Context, key string, ttl time.Duration) error {
	success, fail := c.observe("SET")
	defer fail()

	err := c.client.Set(ctx, key, "", ttl).Err()
	if err != nil {
		fail()
	} else {
		success()
	}

	return err
}

func (c *Cache) Del(ctx context.Context, key string) error {
	success, fail := c.observe("DEL")
	defer fail()

	err := c.client.Del(ctx, key).Err()
	if err != nil {
		fail()
	} else {
		success()
	}

	return err
}

func (c *Cache) observe(method string) (success, fail func()) {
	finish := false
	start := time.Now()

	success = func() {
		if finish {
			return
		}

		finish = true
		duration := time.Since(start).Seconds()

		c.metrics.Requests.WithLabelValues(method, "OK").Inc()
		c.metrics.RequestsDuration.WithLabelValues(method, "OK").Observe(duration)
	}
	fail = func() {
		if finish {
			return
		}

		finish = true
		duration := time.Since(start).Seconds()

		c.metrics.Requests.WithLabelValues(method, "ERR").Inc()
		c.metrics.RequestsDuration.WithLabelValues(method, "ERR").Observe(duration)
	}

	return
}
