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

	q := `
	SELECT o.*, a.street, a.city, a.province, a.postal_code 
		FROM outlets o
	INNER JOIN addresses a 
		ON o.address_id = a.id
	WHERE o.id = $1
	LIMIT 1`

	if err := r.QueryRowxContext(ctx, q, id).StructScan(oa); err != nil {
		return nil, err
	}
	return oa, nil
}

func (r *repository) AddOne(ctx context.Context, no *outlet.NewOutletDTO) (*outlet.SchemaWithAddress, error) {
	o := new(outlet.SchemaWithAddress)

	q := `
	WITH new_address AS (
		INSERT INTO addresses(street, city, province, postal_code)
		VALUES($1, $2, $3, $4)
		RETURNING *
	)
	INSERT INTO outlets(name, phone, opening_time, closing_time, address_id)
	VALUES($5, $6, $7, $8, (SELECT id FROM new_address))
	RETURNING
		*,
		(SELECT street FROM new_address),
		(SELECT city FROM new_address),
		(SELECT province FROM new_address),
		(SELECT postal_code FROM new_address)`

	err := r.QueryRowxContext(ctx, q,
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

	q := `
	UPDATE outlets
	SET name = $1, phone = $2, opening_time = $3, closing_time = $4, updated_at = CURRENT_TIMESTAMP
	WHERE id = $5
	RETURNING *`

	err := r.QueryRowxContext(ctx, q, no.Name, no.Phone, no.OpeningTime, no.ClosingTime, id).
		StructScan(oa)

	if err != nil {
		return nil, err
	}

	return oa, nil
}
