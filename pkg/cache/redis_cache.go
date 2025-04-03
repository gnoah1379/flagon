package cache

import (
	"context"
	"flagon/pkg/config"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	*redis.Client
}

func New() (*RedisCache, error) {
	cfg := config.GetConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Cache.Addr,
		Password: cfg.Cache.Password,
		DB:       cfg.Cache.DB,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &RedisCache{Client: client}, nil
}
