package menutopingrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/goplateframework/internal/domain/menutoping"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	*sqlx.DB
}

func NewDB(db *sqlx.DB) *repository {
	return &repository{db}
}

func (dbrepo *repository) Create(ctx context.Context, m *menutoping.MenuTopingsDTO) error {
	q := `
	INSERT INTO menu_topings
		(id, name, price, is_available, image_url, stock, created_at, updated_at, menu_id)
	VALUES
		(:id, :name, :price, :is_available, :image_url, :stock, :created_at, :updated_at, :menu_id)`

	_, err := dbrepo.NamedExecContext(ctx, q, intoModel(m))
	return err
}

func (dbrepo *repository) GetAll(ctx context.Context) ([]*menutoping.MenuTopingsDTO, error) {
	q := `SELECT * FROM menu_topings`

	var mt []*menutoping.MenuTopingsDTO

	if err := dbrepo.SelectContext(ctx, &mt, q); err != nil {
		return nil, err
	}

	return mt, nil
}

func (dbrepo *repository) GetOne(ctx context.Context, id uuid.UUID) (*menutoping.MenuTopingsDTO, error) {
	mt := new(Model)

	q := `SELECT * FROM menu_topings WHERE id = $1`

	if err := dbrepo.QueryRowxContext(ctx, q, id).StructScan(mt); err != nil {
		return nil, err
	}

	return mt.intoDTO(), nil
}

func (dbrepo *repository) Update(ctx context.Context, m *menutoping.MenuTopingsDTO) error {
	q := `
	UPDATE
		menu_topings
	SET
		name = :name,
		price = :price,
		is_available = :is_available,
		image_url = :image_url,
		stock = :stock,
		updated_at = :updated_at
	WHERE id = :id`

	_, err := dbrepo.NamedExecContext(ctx, q, intoModel(m))
	return err
}

func (dbrepo *repository) Delete(ctx context.Context, id uuid.UUID) error {
	q := `
	DELETE FROM menu_topings
	WHERE id = $1`

	_, err := dbrepo.ExecContext(ctx, q, id)
	return err
}
