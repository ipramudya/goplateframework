package menurepo

import (
	"context"

	"github.com/goplateframework/internal/domain/menu"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	*sqlx.DB
}

func NewDB(db *sqlx.DB) *repository {
	return &repository{db}
}

func (r *repository) AddOne(ctx context.Context, nm *menu.NewMenuDTO) (*menu.Schema, error) {
	m := new(menu.Schema)

	err := r.QueryRowxContext(ctx, createMenuQuery, nm.Name, nm.Description, nm.Price, false, "", nm.OutletID).
		StructScan(m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r *repository) Update(ctx context.Context, nm *menu.NewMenuDTO, id string) (*menu.Schema, error) {
	m := new(menu.Schema)

	err := r.QueryRowxContext(ctx, updateMenuQuery, nm.Name, nm.Description, nm.Price, nm.IsAvailable, "", id).
		StructScan(m)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r *repository) GetAllByOutletID(ctx context.Context, outletID string) (*[]menu.Schema, error) {
	m := []menu.Schema{}

	err := r.SelectContext(ctx, &m, getAllByOutletIDQuery, outletID)
	if err != nil {
		return nil, err
	}

	return &m, nil
}
