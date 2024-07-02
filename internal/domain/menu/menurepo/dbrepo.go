package menurepo

import (
	"context"
	"fmt"
	"strings"

	"github.com/goplateframework/internal/domain/menu"
	"github.com/goplateframework/internal/domain/menu/menuweb"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	*sqlx.DB
}

func NewDB(db *sqlx.DB) *repository {
	return &repository{db}
}

func (dbrepo *repository) Create(ctx context.Context, m *menu.MenuDTO) error {
	query := `
	INSERT INTO menus 
		(id, name, description, price, is_available, image_url, outlet_id, created_at, updated_at)
	VALUES 
		(:id, :name, :description, :price, :is_available, :image_url, :outlet_id, :created_at, :updated_at)`

	_, err := dbrepo.NamedExecContext(ctx, query, intoModel(m))
	return err
}

func (dbrepo *repository) GetAll(ctx context.Context, qp *menuweb.QueryParams) ([]menu.MenuDTO, error) {
	args := map[string]any{
		"size":      qp.Page.Size,
		"outlet_id": qp.Filter.OutletId,
		"offset":    qp.Page.Offset,
	}

	var qb strings.Builder
	qb.WriteString(`
		SELECT * FROM menus
		WHERE outlet_id = :outlet_id
	`)

	if qp.Filter.Name != "" {
		qb.WriteString(" AND name LIKE :name")
		args["name"] = "%" + qp.Filter.Name + "%"
	}

	qb.WriteString(fmt.Sprintf(" ORDER BY %s %s", qp.OrderBy.Field, qp.OrderBy.Direction))
	qb.WriteString(" OFFSET :page LIMIT :size")

	rows, err := dbrepo.NamedQueryContext(ctx, qb.String(), args)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menus []menu.MenuDTO
	for rows.Next() {
		v := new(Model)
		if err := rows.StructScan(v); err != nil {
			return nil, err
		}
		menus = append(menus, *v.intoDTO())
	}

	return menus, nil
}

func (dbrepo *repository) Update(ctx context.Context, nm *menu.MenuDTO) error {
	query := `
	UPDATE 
		menus
	SET 
		name = :name,
		description = :description,
		price = :price,
		is_available = :is_available,
		image_url = :image_url,
		updated_at = :updated_at
	WHERE id = :id`

	_, err := dbrepo.NamedExecContext(ctx, query, intoModel(nm))
	return err
}

func (dbrepo *repository) Count(ctx context.Context, outletId string) (int, error) {
	query := `
	SELECT COUNT(*) AS total FROM menus
	WHERE outlet_id = $1`

	var count struct {
		Total int `db:"total"`
	}

	err := dbrepo.GetContext(ctx, &count, query, outletId)
	if err != nil {
		return 0, err
	}

	return count.Total, nil
}
