package menu

import (
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type QueryParams struct {
	OutletID string `query:"outlet_id"`
	Page     int    `query:"page"`
	Size     int    `query:"size"`
	OrderBy  string `query:"order_by"`
}

const (
	DefaultSize    = 10
	DefaultPage    = 1
	DefaultOrderBy = "id"
)

func ParseQueryParams(c echo.Context) (*QueryParams, error) {
	qp := &QueryParams{
		Page:    DefaultPage,
		Size:    DefaultSize,
		OrderBy: DefaultOrderBy,
	}

	outletIDQuery := c.QueryParam("outlet_id")
	if outletIDQuery == "" {
		return nil, errors.New("outlet_id cannot be empty")
	}

	if outletID, err := uuid.Parse(outletIDQuery); err != nil {
		return nil, err
	} else {
		qp.OutletID = outletID.String()
	}

	if pageQuery := c.QueryParam("page"); pageQuery != "" {
		page, err := strconv.Atoi(pageQuery)
		if err != nil {
			return nil, err
		}
		qp.Page = page
	}

	if sizeQuery := c.QueryParam("size"); sizeQuery != "" {
		size, err := strconv.Atoi(sizeQuery)
		if err != nil {
			return nil, err
		}
		qp.Size = size
	}

	if orderBy := c.QueryParam("order_by"); orderBy != "" {
		qp.OrderBy = orderBy
	}

	return qp, nil
}
