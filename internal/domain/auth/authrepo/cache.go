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

func (c *Cache) AddTokenToBlacklist(ctx context.Context, token string, exp time.Duration) error {
	return c.Set(ctx, token, token, exp).Err()
}
