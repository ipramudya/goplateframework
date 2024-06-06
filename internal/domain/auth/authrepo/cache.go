package authrepo

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	*redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client}
}

func (c *Cache) AddAccessTokenToBlacklist(ctx context.Context, accountID, token string, exp time.Duration) error {
	return c.Set(ctx, accountID, token, exp).Err()
}
