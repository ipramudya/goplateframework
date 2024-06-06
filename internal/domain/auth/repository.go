package auth

import (
	"context"
	"time"
)

type CacheRepository interface {
	AddAccessTokenToBlacklist(ctx context.Context, token string, exp time.Duration) error
	AddRefreshTokenToBlacklist(ctx context.Context, accountID, token string, exp time.Duration) error
	RemoveRefreshTokenFromBlacklist(ctx context.Context, accountID string) error
}
