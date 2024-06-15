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

func (repo repository) Update(ctx context.Context, na *address.NewAddressDTO, id string) (*address.Schema, error) {
	a := new(address.Schema)

	err := repo.QueryRowxContext(ctx, updateQuery, na.Street, na.City, na.Province, na.PostalCode, id).
		StructScan(a)

	if err != nil {
		return &address.Schema{}, err
	}

	return a, nil
}

func (repo repository) GetOneByID(ctx context.Context, id string) (*address.Schema, error) {
	a := new(address.Schema)

	err := repo.QueryRowxContext(ctx, getOneByIDQuery, id).StructScan(a)

	if err != nil {
		return &address.Schema{}, err
	}

	return a, nil
}
