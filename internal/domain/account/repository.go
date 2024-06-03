package account

import (
	"context"
	"time"
)

type DBRepository interface {
	GetOneByEmail(ctx context.Context, email string) (*Schema, error)
	Register(ctx context.Context, account *NewAccouuntDTO) (*Schema, error)
	ChangePassword(ctx context.Context, email, password string) error
}

type CacheRepository interface {
	AddTokenToBlacklist(ctx context.Context, token string, exp time.Duration) error
}
