package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type TokenRepository interface {
	AddUserToken(ctx context.Context, userID uuid.UUID, jti string, expiration time.Duration) error
	RemoveUserToken(ctx context.Context, userID uuid.UUID, jti string) error
	IsUserTokenValid(ctx context.Context, userID uuid.UUID, jti string) (bool, error)
}

type redisTokenRepo struct {
	client *redis.Client
}

func NewRedisTokenRepository(client *redis.Client) TokenRepository {
	return &redisTokenRepo{
		client: client,
	}
}

func (r *redisTokenRepo) AddUserToken(ctx context.Context, userID uuid.UUID, jti string, expiration time.Duration) error {
	key := fmt.Sprintf("user:%s:tokens:%s", userID.String(), jti)
	return r.client.Set(ctx, key, jti, expiration).Err()
}

func (r *redisTokenRepo) RemoveUserToken(ctx context.Context, userID uuid.UUID, jti string) error {
	key := fmt.Sprintf("user:%s:tokens:%s", userID.String(), jti)

	return r.client.Del(ctx, key).Err()
}

func (r *redisTokenRepo) IsUserTokenValid(ctx context.Context, userID uuid.UUID, jti string) (bool, error) {
	key := fmt.Sprintf("user:%s:tokens:%s", userID.String(), jti)

	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}
