package authrepo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	*redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client}
}

// stored access token using token itself as key on redis
func (c *Cache) AddAccessTokenToBlacklist(ctx context.Context, token string, exp time.Duration) error {
	return c.Set(ctx, token, token, exp).Err()
}

// stored refresh token using account id as key on redis
func (c *Cache) AddRefreshTokenToBlacklist(ctx context.Context, accountID uuid.UUID, token string, exp time.Duration) error {
	return c.Set(ctx, accountID.String(), token, exp).Err()
}

// remove existing refresh token from redis using account id to make redis clean
func (c *Cache) RemoveRefreshTokenFromBlacklist(ctx context.Context, accountID uuid.UUID) error {
	_, err := c.Get(ctx, accountID.String()).Result()

	if err == redis.Nil {
		return nil
	}

	return c.Del(ctx, accountID.String()).Err()
}
