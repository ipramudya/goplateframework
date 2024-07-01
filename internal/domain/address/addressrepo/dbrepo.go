package addressrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/goplateframework/internal/domain/address"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	*sqlx.DB
}

func NewDB(db *sqlx.DB) *repository {
	return &repository{db}
}

func (dbrepo *repository) Create(ctx context.Context, a *address.AddressDTO) error {
	query := `
	INSERT INTO addresses 
		(id, street, city, province, postal_code, created_at, updated_at)
	VALUES
		(:id, :street, :city, :province, :postal_code, :created_at, :updated_at)
	`

	_, err := dbrepo.NamedExecContext(ctx, query, intoModel(a))
	return err
}

func (dbrepo *repository) GetOne(ctx context.Context, id uuid.UUID) (*address.AddressDTO, error) {
	query := `
	SELECT * FROM addresses
	WHERE id = $1
	LIMIT 1
	`

	a := new(address.AddressDTO)
	err := dbrepo.QueryRowContext(ctx, query, id).Scan(a)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (dbrepo *repository) Update(ctx context.Context, na *address.AddressDTO) error {
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

func (dbrepo *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE FROM addresses
	WHERE id = $1
	`

	_, err := dbrepo.ExecContext(ctx, query, id)
	return err
}
