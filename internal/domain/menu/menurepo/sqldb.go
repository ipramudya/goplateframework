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

	q := `
	INSERT INTO menus (name, description, price, is_available, image_url, outlet_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING *`

	err := r.QueryRowxContext(ctx, q, nm.Name, nm.Description, nm.Price, false, "", nm.OutletID).
		StructScan(m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r *repository) Update(ctx context.Context, nm *menu.NewMenuDTO, id string) (*menu.Schema, error) {
	m := new(menu.Schema)

	q := `
	UPDATE menus
	SET name = $1, description = $2, price = $3, is_available = $4, image_url = $5, updated_at = CURRENT_TIMESTAMP
	WHERE id = $6
	RETURNING *`

	err := r.QueryRowxContext(ctx, q, nm.Name, nm.Description, nm.Price, nm.IsAvailable, "", id).
		StructScan(m)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r *repository) GetAllByOutletID(ctx context.Context, outletID string) (*[]menu.Schema, error) {
	m := []menu.Schema{}

	q := `
	SELECT * FROM menus
	WHERE outlet_id = $1
	ORDER BY created_at DESC`

	err := r.SelectContext(ctx, &m, q, outletID)
	if err != nil {
		return nil, err
	}

	return &m, nil
}
