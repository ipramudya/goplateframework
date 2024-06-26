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

	q := `
	SELECT * FROM accounts 
	WHERE email=$1
	LIMIT 1`

	if err := repo.QueryRowxContext(ctx, q, email).StructScan(account); err != nil {
		return nil, err
	}

	return account, nil
}

func (repo repository) GetOneByID(ctx context.Context, id string) (*account.Schema, error) {
	account := new(account.Schema)

	q := `
	SELECT * FROM accounts 
	WHERE id=$1
	LIMIT 1`

	if err := repo.QueryRowxContext(ctx, q, id).StructScan(account); err != nil {
		return nil, err
	}

	return account, nil
}

func (repo repository) Register(ctx context.Context, na *account.NewAccouuntDTO) (*account.Schema, error) {
	account := new(account.Schema)

	q := `
	INSERT INTO accounts(firstname, lastname, email, password, phone)
	VALUES($1, $2, $3, $4, $5)
	RETURNING *`

	err := repo.
		QueryRowxContext(ctx, q, na.Firstname, na.Lastname, na.Email, na.Password, na.Phone).
		StructScan(account)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (repo repository) ChangePassword(ctx context.Context, email, password string) error {
	q := `
	UPDATE accounts
	SET password=$1, updated_at = CURRENT_TIMESTAMP
	WHERE email=$2`

	_, err := repo.ExecContext(ctx, q, password, email)
	return err
}
