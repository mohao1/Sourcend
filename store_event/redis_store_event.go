package store_event

import (
	"context"
	"crypto/tls"
	"github.com/redis/go-redis/v9"
)

type RedisStoreConfig struct {
	Addr      string
	Pass      string
	DB        int
	IsTLS     bool
	TLSConfig *tls.Config
	IsLimiter bool
	Limiter   redis.Limiter
}

// RedisStore Redis的StoreEvent的实现
type RedisStore struct {
	redisClient *redis.Client
	config      RedisStoreConfig
}

func NewRedisStore(config RedisStoreConfig) *RedisStore {
	redisClient := redis.NewClient(&redis.Options{
		Addr:      config.Addr,
		Password:  config.Pass,
		DB:        config.DB,
		TLSConfig: config.TLSConfig,
		Limiter:   config.Limiter,
	})

	return &RedisStore{
		redisClient: redisClient,
		config:      config,
	}
}

func (r *RedisStore) Handler(ctx context.Context, data StoreEventInfo) error {
	//action.NewAction()
	//r.redisClient.HSet()
	return nil
}
