package cache

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-redis/redismock/v9"

	"github.com/AydarZaynutdinov/pow/internal/metrics"
)

const (
	key = "mock_key"
	ttl = time.Second
)

var (
	cacheMetrics = metrics.NewCacheMetrics()

	errCustom = errors.New("custom error")
)

func TestCache_Del(t *testing.T) {
	tests := []struct {
		name         string
		prepareCache func() *Cache
		check        func(err error)
	}{
		{
			name: "error during del",
			prepareCache: func() *Cache {
				client, mock := redismock.NewClientMock()

				cache := &Cache{
					client:  client,
					metrics: cacheMetrics,
				}

				mock.ExpectDel(key).SetErr(errCustom)

				return cache
			},
			check: func(err error) {
				assert.Equal(t, errCustom, err)
			},
		},
		{
			name: "without error during del",
			prepareCache: func() *Cache {
				client, mock := redismock.NewClientMock()

				cache := &Cache{
					client:  client,
					metrics: cacheMetrics,
				}

				mock.ExpectDel(key).SetVal(1)

				return cache
			},
			check: func(err error) {
				assert.Nil(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := tt.prepareCache()
			err := cache.Del(context.Background(), key)
			tt.check(err)
		})
	}
}

func TestCache_Exists(t *testing.T) {
	tests := []struct {
		name         string
		prepareCache func() *Cache
		check        func(res int64, err error)
	}{
		{
			name: "error during exists",
			prepareCache: func() *Cache {
				client, mock := redismock.NewClientMock()

				cache := &Cache{
					client:  client,
					metrics: cacheMetrics,
				}

				mock.ExpectExists(key).SetErr(errCustom)

				return cache
			},
			check: func(res int64, err error) {
				assert.Equal(t, int64(0), res)
				assert.Equal(t, errCustom, err)
			},
		},
		{
			name: "without error during exists | return 0",
			prepareCache: func() *Cache {
				client, mock := redismock.NewClientMock()

				cache := &Cache{
					client:  client,
					metrics: cacheMetrics,
				}

				mock.ExpectExists(key).SetVal(0)

				return cache
			},
			check: func(res int64, err error) {
				assert.Equal(t, int64(0), res)
				assert.Nil(t, err)
			},
		},
		{
			name: "without error during exists | return 1",
			prepareCache: func() *Cache {
				client, mock := redismock.NewClientMock()

				cache := &Cache{
					client:  client,
					metrics: cacheMetrics,
				}

				mock.ExpectExists(key).SetVal(1)

				return cache
			},
			check: func(res int64, err error) {
				assert.Equal(t, int64(1), res)
				assert.Nil(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := tt.prepareCache()
			res, err := cache.Exists(context.Background(), key)
			tt.check(res, err)
		})
	}
}

func TestCache_Set(t *testing.T) {
	tests := []struct {
		name         string
		prepareCache func() *Cache
		check        func(err error)
	}{
		{
			name: "error during set",
			prepareCache: func() *Cache {
				client, mock := redismock.NewClientMock()

				cache := &Cache{
					client:  client,
					metrics: cacheMetrics,
				}

				mock.ExpectSet(key, "", ttl).SetErr(errCustom)

				return cache
			},
			check: func(err error) {
				assert.Equal(t, errCustom, err)
			},
		},
		{
			name: "without error during set",
			prepareCache: func() *Cache {
				client, mock := redismock.NewClientMock()

				cache := &Cache{
					client:  client,
					metrics: cacheMetrics,
				}

				mock.ExpectSet(key, "", ttl).SetVal("")

				return cache
			},
			check: func(err error) {
				assert.Nil(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := tt.prepareCache()
			err := cache.Set(context.Background(), key, ttl)
			tt.check(err)
		})
	}
}
