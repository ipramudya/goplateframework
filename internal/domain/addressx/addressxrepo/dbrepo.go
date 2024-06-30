package addressxrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/goplateframework/internal/domain/addressx"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	*sqlx.DB
}

func NewDB(db *sqlx.DB) *repository {
	return &repository{db}
}

func (dbrepo repository) Create(ctx context.Context, a *addressx.AddressDTO) error {
	query := `
	INSERT INTO addresses 
		(id, street, city, province, postal_code, created_at, updated_at)
	VALUES
		(:id, :street, :city, :province, :postal_code, :created_at, :updated_at)
	`

	_, err := dbrepo.NamedExecContext(ctx, query, intoModel(a))
	return err
}

func (dbrepo repository) GetOne(ctx context.Context, id uuid.UUID) (*addressx.AddressDTO, error) {
	query := `
	SELECT * FROM addresses
	WHERE id = $1
	LIMIT 1
	`

	a := new(addressx.AddressDTO)
	err := dbrepo.QueryRowContext(ctx, query, id).Scan(a)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (dbrepo repository) Update(ctx context.Context, na *addressx.AddressDTO) error {
	query := `
	UPDATE
		addresses
	SET
		street = :street,
		city = :city,
		province = :province,
		postal_code = :postal_code
		updated_at = :updated_at
	WHERE id = :id
	`

	_, err := dbrepo.NamedExecContext(ctx, query, intoModel(na))
	if err != nil {
		return err
	}

	return nil
}

func (dbrepo repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE FROM addresses
	WHERE id = $1
	`

	_, err := dbrepo.ExecContext(ctx, query, id)
	return err
}
