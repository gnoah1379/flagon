package repository

import (
	"context"
	"flagon/pkg/cache"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TokenRepository interface {
	AddJwtToken(ctx context.Context, userID uuid.UUID, jti string, expiration time.Duration) error
	RemoveJwtToken(ctx context.Context, userID uuid.UUID, jti string) error
	IsJwtValid(ctx context.Context, userID uuid.UUID, jti string) (bool, error)
}

type redisTokenRepo struct {
	client *cache.RedisCache
}

func NewTokenRepository(client *cache.RedisCache) TokenRepository {
	return &redisTokenRepo{
		client: client,
	}
}

func (r *redisTokenRepo) AddJwtToken(ctx context.Context, userID uuid.UUID, jti string, expiration time.Duration) error {
	key := fmt.Sprintf("user:%s:jwt-tokens:%s", userID.String(), jti)
	return r.client.Set(ctx, key, jti, expiration).Err()
}

func (r *redisTokenRepo) RemoveJwtToken(ctx context.Context, userID uuid.UUID, jti string) error {
	key := fmt.Sprintf("user:%s:jwt-tokens:%s", userID.String(), jti)

	return r.client.Del(ctx, key).Err()
}

func (r *redisTokenRepo) IsJwtValid(ctx context.Context, userID uuid.UUID, jti string) (bool, error) {
	key := fmt.Sprintf("user:%s:jwt-tokens:%s", userID.String(), jti)

	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}
