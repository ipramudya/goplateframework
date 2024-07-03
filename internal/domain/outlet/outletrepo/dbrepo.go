package outletrepo

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/goplateframework/internal/domain/outlet"
	"github.com/goplateframework/internal/domain/outlet/outletweb"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	*sqlx.DB
}

func NewDB(db *sqlx.DB) *repository {
	return &repository{db}
}

func (dbrepo *repository) GetOne(ctx context.Context, id uuid.UUID) (*outlet.OutletDTO, error) {
	oa := new(ModelWithAddress)

	q := `
	SELECT o.*, a.street, a.city, a.province, a.postal_code 
		FROM outlets o
	INNER JOIN addresses a 
		ON o.address_id = a.id
	WHERE o.id = $1
	LIMIT 1`

	if err := dbrepo.QueryRowxContext(ctx, q, id).StructScan(oa); err != nil {
		return nil, err
	}

	return oa.intoDTO(), nil
}

func (dbrepo *repository) Create(ctx context.Context, o *outlet.OutletDTO) error {
	q := `
	INSERT INTO outlets
		(name, phone, opening_time, closing_time, address_id, created_at, updated_at)
	VALUES
		(:name, :phone, :opening_time, :closing_time, :address_id, :created_at, :updated_at)`

	_, err := dbrepo.NamedExecContext(ctx, q, intoModel(o))
	return err
}

func (dbrepo *repository) Update(ctx context.Context, o *outlet.OutletDTO) error {
	q := `
	UPDATE 
		outlets
	SET
		name = :name,
		phone = :phone,
		opening_time = :opening_time,
		closing_time = :closing_time,
		updated_at = :updated_at
	WHERE id = :id`

	_, err := dbrepo.NamedExecContext(ctx, q, intoModel(o))
	return err
}

func (dbrepo *repository) Count(ctx context.Context) (int, error) {
	q := `
	SELECT COUNT(*) AS total FROM outlets`

	var count struct {
		Total int `db:"total"`
	}

	err := dbrepo.GetContext(ctx, &count, q)
	if err != nil {
		return 0, err
	}

	return count.Total, nil
}

func (dbrepo *repository) GetAll(ctx context.Context, qp *outletweb.QueryParams) ([]outlet.OutletDTO, error) {
	args := map[string]any{
		"size":   qp.Page.Size,
		"offset": qp.Page.Offset,
	}

	var qb strings.Builder
	qb.WriteString(`
		SELECT 
			o.*, a.street, a.city, a.province, a.postal_code
		FROM 
			outlets o
		INNER JOIN addresses a
			ON o.address_id = a.id
	`)

	qb.WriteString(dbrepo.buildFilter(args, qp))
	qb.WriteString(fmt.Sprintf(" ORDER BY o.%s %s", qp.OrderBy.Field, qp.OrderBy.Direction))
	qb.WriteString(" OFFSET :offset LIMIT :size")

	rows, err := dbrepo.NamedQueryContext(ctx, qb.String(), args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var outlets []outlet.OutletDTO
	for rows.Next() {
		o := new(ModelWithAddress)
		if err := rows.StructScan(o); err != nil {
			return nil, err
		}
		outlets = append(outlets, *o.intoDTO())
	}

	return outlets, nil
}

func (dbrepo *repository) Delete(ctx context.Context, id uuid.UUID) error {
	q := `
	DELETE FROM outlets
	WHERE id = $1`

	_, err := dbrepo.ExecContext(ctx, q, id)
	return err
}
