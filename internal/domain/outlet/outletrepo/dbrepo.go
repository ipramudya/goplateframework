package outletrepo

import (
	"bytes"
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

func (dbrepo *repository) Create(ctx context.Context, o *outlet.OutletDTO) error {
	query := `
	INSERT INTO outlets
		(name, phone, opening_time, closing_time, address_id, created_at, updated_at)
	VALUES
		(:name, :phone, :opening_time, :closing_time, :address_id, :created_at, :updated_at)`

	_, err := dbrepo.NamedExecContext(ctx, query, intoModel(o))
	return err
}

func (dbrepo *repository) Update(ctx context.Context, o *outlet.OutletDTO) error {
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

func (dbrepo *repository) Count(ctx context.Context) (int, error) {
	query := `
	SELECT COUNT(*) AS total FROM outlets`

	var count struct {
		Total int `db:"total"`
	}

	err := dbrepo.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}

	return count.Total, nil
}

func (dbrepo *repository) GetAll(ctx context.Context, qp *outletweb.QueryParams) ([]outlet.OutletDTO, error) {
	args := map[string]any{
		"last_id": qp.Filter.LastId,
		"size":    qp.Page.Size,
	}

	var query string

	if args["last_id"] == "" {
		query = `
		SELECT
			o.*, a.street, a.city, a.province, a.postal_code 
		FROM 
			outlets o
		INNER JOIN addresses a 
			ON o.address_id = a.id
		{filter}
		{order_by}
		LIMIT :size`
	} else {
		query = `
		SELECT * FROM (
			SELECT 
				o.*, a.street, a.city, a.province, a.postal_code 
			FROM 
				outlets o
			INNER JOIN addresses a
				ON o.address_id = a.id
			WHERE 
				o.created_at < (
					SELECT created_at FROM outlets
					WHERE id = :last_id
				)
			ORDER BY o.created_at DESC
			LIMIT :size
		) AS o
		{filter}
		{order_by}`
	}

	queryByte := []byte(query)

	filterByte := dbrepo.useFilter(args, qp)
	queryByte = bytes.Replace(queryByte, []byte("{filter}"), filterByte, -1)

	orderbyByte := []byte(fmt.Sprintf(" ORDER BY o.%s %s", qp.OrderBy.Field, qp.OrderBy.Direction))
	queryByte = bytes.Replace(queryByte, []byte("{order_by}"), orderbyByte, -1)

	query = strings.Join(strings.Fields(string(queryByte)), " ")

	rows, err := dbrepo.NamedQueryContext(ctx, query, args)
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
		fmt.Printf("MODEL : %+v\n", o)
		outlets = append(outlets, *o.intoDTO())
	}

	return outlets, nil
}
