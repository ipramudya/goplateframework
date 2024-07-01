package accountrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/goplateframework/internal/domain/account"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	*sqlx.DB
}

func NewDB(db *sqlx.DB) *repository {
	return &repository{db}
}

func (dbrepo *repository) GetOneByEmail(ctx context.Context, email string) (*account.AccountDTO, error) {
	a := new(Model)

	query := `
	SELECT * FROM accounts
	WHERE email = $1
	LIMIT 1
	`

	if err := dbrepo.QueryRowxContext(ctx, query, email).StructScan(a); err != nil {
		return nil, err
	}

	return a.intoDTO(), nil
}

func (dbrepo *repository) GetOne(ctx context.Context, id uuid.UUID) (*account.AccountDTO, error) {
	a := new(Model)

	query := `
	SELECT * FROM accounts
	WHERE id = $1
	LIMIT 1
	`

	if err := dbrepo.QueryRowxContext(ctx, query, id).StructScan(a); err != nil {
		return nil, err
	}

	return a.intoDTO(), nil
}

func (dbrepo *repository) Create(ctx context.Context, a *account.AccountDTO) error {
	query := `
	INSERT INTO accounts
		(id, firstname, lastname, email, password, phone, role, created_at, updated_at)
	VALUES
		(:id, :firstname, :lastname, :email, :password, :phone, :role, :created_at, :updated_at)`

	_, err := dbrepo.NamedExecContext(ctx, query, intoModel(a))
	return err
}

func (repo *repository) ChangePassword(ctx context.Context, email, password string) error {
	query := `
	UPDATE accounts
	SET password = $1
	WHERE email = $2`

	_, err := repo.ExecContext(ctx, query, password, email)
	return err
}
