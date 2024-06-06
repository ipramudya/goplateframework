package auth

import (
	"context"
	"time"
)

type CacheRepository interface {
	AddAccessTokenToBlacklist(ctx context.Context, accountID, token string, exp time.Duration) error
}
