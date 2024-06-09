package account

import (
	"context"
)

type DBRepository interface {
	GetOneByEmail(ctx context.Context, email string) (*Schema, error)
	GetOneByID(ctx context.Context, id string) (*Schema, error)
	Register(ctx context.Context, account *NewAccouuntDTO) (*Schema, error)
	ChangePassword(ctx context.Context, email, password string) error
}
