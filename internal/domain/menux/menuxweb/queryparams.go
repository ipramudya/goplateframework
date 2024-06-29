package menuxweb

import (
	"errors"

	"github.com/google/uuid"
	"github.com/goplateframework/internal/web/queryparams"
	"github.com/labstack/echo/v4"
)

// Supported query params for this menu web layer
type UnparsedQueryParams struct {
	page     string
	size     string
	orderBy  string
	outletId string
	lastId   string
	name     string
}

func getQueryParams(c echo.Context) *UnparsedQueryParams {
	return &UnparsedQueryParams{
		page:     c.QueryParam("page"),
		size:     c.QueryParam("size"),
		orderBy:  c.QueryParam("order_by"),
		outletId: c.QueryParam("outlet_id"),
		lastId:   c.QueryParam("last_id"),
		name:     c.QueryParam("name"),
	}
}

// Populated query params to send to repository
type QueryParams struct {
	Page    *queryparams.Page
	OrderBy *queryparams.OrderBy
	Filter  struct {
		OutletId string
		LastId   string
		Name     string
	}
}

func (uqp *UnparsedQueryParams) Parse() (*QueryParams, error) {
	qp := new(QueryParams)

	if err := uqp.setPage(qp); err != nil {
		return nil, err
	}

	if err := uqp.setOrderBy(qp); err != nil {
		return nil, err
	}

	if err := uqp.setFilter(qp); err != nil {
		return nil, err
	}

	return qp, nil
}

func (uqp *UnparsedQueryParams) setPage(qp *QueryParams) error {
	page, err := queryparams.ParsePage(uqp.page, uqp.size)
	if err != nil {
		return err
	}

	qp.Page = page
	return nil
}

var allowedOrderByFields = []string{"name", "price", "created_at"}

func (uqp *UnparsedQueryParams) setOrderBy(qp *QueryParams) error {
	defaultOrderBy := queryparams.NewOrderBy(
		"updated_at",
		queryparams.AscOrder,
	)

	orderBy, err := queryparams.ParseOrderBy(allowedOrderByFields, uqp.orderBy, defaultOrderBy)
	if err != nil {
		return err
	}

	qp.OrderBy = orderBy
	return nil
}

func (uqp *UnparsedQueryParams) setFilter(qp *QueryParams) error {
	// outlet_id is needed to retrieve menus only for an outlet,
	// system doesn't allow to retrieve all menus
	if uqp.outletId == "" {
		return errors.New("outlet_id cannot be empty")
	}

	outletId, err := uuid.Parse(uqp.outletId)
	if err != nil {
		return errors.New("outlet_id is not valid")
	}
	qp.Filter.OutletId = outletId.String()

	if uqp.lastId != "" {
		if uqp.page == "" && uqp.size == "" {
			return errors.New("page and size cannot be empty when last_id is set")
		}

		if lastId, err := uuid.Parse(uqp.lastId); err != nil {
			return err
		} else {
			qp.Filter.LastId = lastId.String()
		}
	} else {
		qp.Filter.LastId = ""
	}

	qp.Filter.Name = uqp.name

	return nil
}
