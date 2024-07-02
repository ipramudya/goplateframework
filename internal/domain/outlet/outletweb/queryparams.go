package outletweb

import (
	"errors"
	"slices"

	"github.com/goplateframework/internal/web/queryparams"
	"github.com/labstack/echo/v4"
)

// Supported query params for this outlet web layer
type UnparsedQueryParams struct {
	page    string
	size    string
	orderBy string
	name    string
	operate string // open | close
}

func getQueryParams(c echo.Context) *UnparsedQueryParams {
	return &UnparsedQueryParams{
		page:    c.QueryParam("page"),
		size:    c.QueryParam("size"),
		orderBy: c.QueryParam("order_by"),
		name:    c.QueryParam("name"),
		operate: c.QueryParam("operate"),
	}
}

// Populated query params to send to repository
type QueryParams struct {
	Page    *queryparams.Page
	OrderBy *queryparams.OrderBy
	Filter  struct {
		Name    string
		Operate string
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

var allowedOrderByFields = []string{"name", "opening_time", "created_at"}

func (uqp *UnparsedQueryParams) setOrderBy(qp *QueryParams) error {
	defaultOrderBy := queryparams.NewOrderBy(
		"created_at",
		queryparams.DescOrder,
	)

	orderBy, err := queryparams.ParseOrderBy(allowedOrderByFields, uqp.orderBy, defaultOrderBy)
	if err != nil {
		return err
	}

	qp.OrderBy = orderBy
	return nil
}

var operateEnums = []string{"open", "close"}

func (uqp *UnparsedQueryParams) setFilter(qp *QueryParams) error {
	if uqp.operate != "" && !slices.Contains(operateEnums, uqp.operate) {
		return errors.New("operate must be either 'close' or 'open")
	}

	qp.Filter.Operate = uqp.operate
	qp.Filter.Name = uqp.name

	return nil
}
