package account

import "context"

// AccountRepo interface declares the behavior this package needs to perists and
// retrieve data.
type Repository interface {
	GetOneByEmail(ctx context.Context, email string) (*Schema, error)
	Register(ctx context.Context, account *NewAccouuntDTO) (*Schema, error)
	ChangePassword(ctx context.Context, email, password string) error
}
