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
	account := new(account.Schema)

	if err := repo.QueryRowxContext(ctx, getOneByEmailQuery, email).StructScan(account); err != nil {
		return nil, err
	}

	return account, nil
}

func (repo repository) GetOneByID(ctx context.Context, id string) (*account.Schema, error) {
	account := new(account.Schema)

	if err := repo.QueryRowxContext(ctx, getOneByIDQuery, id).StructScan(account); err != nil {
		return nil, err
	}

	return account, nil
}

func (repo repository) Register(ctx context.Context, na *account.NewAccouuntDTO) (*account.Schema, error) {
	account := new(account.Schema)

	err := repo.
		QueryRowxContext(ctx, createAccountQuery, na.Firstname, na.Lastname, na.Email, na.Password, na.Phone).
		StructScan(account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (repo repository) ChangePassword(ctx context.Context, email, password string) error {
	_, err := repo.ExecContext(ctx, changePasswordQuery, password, email)
	return err
}
