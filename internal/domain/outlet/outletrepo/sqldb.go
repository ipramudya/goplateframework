package outletrepo

import (
	"context"

	"github.com/goplateframework/internal/domain/outlet"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	*sqlx.DB
}

func NewDB(db *sqlx.DB) *repository {
	return &repository{db}
}

func (r *repository) GetOneByID(ctx context.Context, id string) (*outlet.SchemaWithAddress, error) {
	oa := new(outlet.SchemaWithAddress)
	if err := r.QueryRowxContext(ctx, getOneByIDQuery, id).StructScan(oa); err != nil {
		return nil, err
	}
	return oa, nil
}

func (r *repository) AddOne(ctx context.Context, no *outlet.NewOutletDTO) (*outlet.SchemaWithAddress, error) {
	o := new(outlet.SchemaWithAddress)

	err := r.QueryRowxContext(ctx, createOutletQuery,
		no.Address.Street, no.Address.City, no.Address.Province, no.Address.PostalCode,
		no.Name, no.Phone, no.OpeningTime, no.ClosingTime).
		StructScan(o)

	if err != nil {
		return nil, err
	}

	return o, nil
}

func (r *repository) Update(ctx context.Context, no *outlet.NewOutletDTO, id string) (*outlet.Schema, error) {
	oa := new(outlet.Schema)

	err := r.QueryRowxContext(ctx, updateOutletQuery, no.Name, no.Phone, no.OpeningTime, no.ClosingTime, id).
		StructScan(oa)

	if err != nil {
		return nil, err
	}

	return oa, nil
}
