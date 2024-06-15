package accountrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/goplateframework/internal/domain/account"
	"github.com/redis/go-redis/v9"
)

const (
	meExpires = 24 * time.Hour
)

type Cache struct {
	*redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client}
}

func (c *Cache) SetMe(ctx context.Context, accountPayload interface{}) error {
	data, err := sonic.Marshal(accountPayload)
	if err != nil {
		return err
	}

	key := keyofMe(accountPayload.(*account.Schema).ID.String())
	return c.Set(ctx, key, data, meExpires).Err()
}

func (c *Cache) GetMe(ctx context.Context, id string) (*account.Schema, error) {
	key := keyofMe(id)

	data, err := c.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	a := new(account.Schema)
	if err := sonic.Unmarshal([]byte(data), &a); err != nil {
		return nil, err
	}

	return a, nil
}

func (c *Cache) RemoveMe(ctx context.Context, id string) error {
	return c.Del(ctx, keyofMe(id)).Err()
}

func keyofMe(id string) string {
	return fmt.Sprintf("me:%s", id)
}
