package accountrepo

import (
	"context"

	"github.com/goplateframework/internal/domain/account"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	*sqlx.DB
}

func NewDB(db *sqlx.DB) *repository {
	return &repository{db}
}

func (repo repository) GetOneByEmail(ctx context.Context, email string) (*account.Schema, error) {
	account := &account.Schema{}

	err := repo.
		QueryRowxContext(ctx, GetOneByEmailQuery, email).
		StructScan(account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (repo repository) GetOneByID(ctx context.Context, id string) (*account.Schema, error) {
	account := &account.Schema{}

	err := repo.
		QueryRowxContext(ctx, GetOneByIDQuery, id).
		StructScan(account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (repo repository) Register(ctx context.Context, na *account.NewAccouuntDTO) (*account.Schema, error) {
	account := &account.Schema{}

	err := repo.
		QueryRowxContext(ctx, CreateAccountQuery, na.Firstname, na.Lastname, na.Email, na.Password, na.Phone).
		StructScan(account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (repo repository) ChangePassword(ctx context.Context, email, password string) error {
	_, err := repo.ExecContext(ctx, ChangePasswordQuery, password, email)
	return err
}
