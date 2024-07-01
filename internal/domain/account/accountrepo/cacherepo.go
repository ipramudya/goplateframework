package accountrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"
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

func (c *Cache) SetMe(ctx context.Context, accountPayload *account.AccountDTO) error {
	data, err := sonic.Marshal(accountPayload)

	if err != nil {
		return err
	}

	key := getMeKey(accountPayload.ID)
	return c.Set(ctx, key, data, meExpires).Err()
}

func (c *Cache) GetMe(ctx context.Context, id uuid.UUID) (*account.AccountDTO, error) {
	key := getMeKey(id)

	data, err := c.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	a := new(Model)
	if err := sonic.Unmarshal([]byte(data), &a); err != nil {
		return nil, err
	}

	return a.intoDTO(), nil
}

func (c *Cache) RemoveMe(ctx context.Context, id uuid.UUID) error {
	return c.Del(ctx, getMeKey(id)).Err()
}

func getMeKey(id uuid.UUID) string {
	return fmt.Sprintf("me:%s", id.String())
}
