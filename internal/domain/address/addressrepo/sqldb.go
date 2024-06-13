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
	address := new(address.Schema)
	if err := repo.QueryRowxContext(ctx, getOneByIDQuery, id).StructScan(address); err != nil {
		return nil, err
	}
	return address, nil
}

func (repo repository) AddOne(ctx context.Context, a *address.NewAddressDTO) (*address.Schema, error) {
	address := new(address.Schema)

	err := repo.
		QueryRowxContext(ctx, addOneQuery, a.Street, a.City, a.Province, a.PostalCode).
		StructScan(address)

	if err != nil {
		return nil, err
	}

	return address, nil
}

func (repo repository) Update(ctx context.Context, a *address.NewAddressDTO) error {
	_, err := repo.ExecContext(ctx, updateQuery, a.Street, a.City, a.Province, a.PostalCode, a.ID)
	return err
}
