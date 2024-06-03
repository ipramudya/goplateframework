package auth

import (
	"context"
	"time"
)

type CacheRepository interface {
	AddTokenToBlacklist(ctx context.Context, token string, exp time.Duration) error
}
