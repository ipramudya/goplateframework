package outletxrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/goplateframework/internal/domain/outletx"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	*sqlx.DB
}

func NewDB(db *sqlx.DB) *repository {
	return &repository{db}
}

func (dbrepo repository) GetOne(ctx context.Context, id uuid.UUID) (*outletx.OutletDTO, error) {
	oa := new(ModelWithAddress)

	query := `
	SELECT o.*, a.street, a.city, a.province, a.postal_code 
		FROM outlets o
	INNER JOIN addresses a 
		ON o.address_id = a.id
	WHERE o.id = $1
	LIMIT 1`

	if err := dbrepo.QueryRowxContext(ctx, query, id).StructScan(oa); err != nil {
		return nil, err
	}

	return oa.intoDTO(), nil
}

func (dbrepo repository) Create(ctx context.Context, o *outletx.OutletDTO) error {
	query := `
	INSERT INTO outlets
		(name, phone, opening_time, closing_time, address_id, created_at, updated_at)
	VALUES
		(:name, :phone, :opening_time, :closing_time, :address_id, :created_at, :updated_at)`

	_, err := dbrepo.NamedExecContext(ctx, query, intoModel(o))
	return err
}

func (dbrepo repository) Update(ctx context.Context, o *outletx.OutletDTO) error {
	query := `
	UPDATE 
		outlets
	SET
		name = :name,
		phone = :phone,
		opening_time = :opening_time,
		closing_time = :closing_time,
		updated_at = :updated_at
	WHERE id = :id`

	_, err := dbrepo.NamedExecContext(ctx, query, intoModel(o))
	return err
}
