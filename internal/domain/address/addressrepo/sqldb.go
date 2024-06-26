package addressrepo

import (
	"context"

	"github.com/goplateframework/internal/domain/address"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	*sqlx.DB
}

func NewDB(db *sqlx.DB) *repository {
	return &repository{db}
}

func (repo repository) GetOneByID(ctx context.Context, id string) (*address.Schema, error) {
	a := new(address.Schema)

	q := `
	SELECT * FROM addresses
	WHERE id = $1
	LIMIT 1`

	err := repo.QueryRowxContext(ctx, q, id).StructScan(a)

	if err != nil {
		return nil, err
	}

	return a, nil
}

func (repo repository) Update(ctx context.Context, na *address.NewAddressDTO, id string) (*address.Schema, error) {
	a := new(address.Schema)

	q := `
	UPDATE addresses
	SET street = $1, city = $2, province = $3, postal_code = $4, updated_at = CURRENT_TIMESTAMP
	WHERE id = $5
	RETURNING *`

	err := repo.QueryRowxContext(ctx, q, na.Street, na.City, na.Province, na.PostalCode, id).
		StructScan(a)

	if err != nil {
		return nil, err
	}

	return a, nil
}
