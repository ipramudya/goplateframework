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

func (repo repository) Update(ctx context.Context, a *address.AddressDTO) error {
	_, err := repo.ExecContext(ctx, updateQuery, a.Street, a.City, a.Province, a.PostalCode, a.ID)
	return err
}
