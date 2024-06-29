package queryparams

import (
	"fmt"
	"slices"
	"strings"
)

const (
	AscOrder  = "ASC"
	DescOrder = "DESC"
)

var directions = []string{AscOrder, DescOrder}

type OrderBy struct {
	Direction string
	Field     string
}

func NewOrderBy(field, direction string) *OrderBy {
	if !slices.Contains(directions, direction) {
		return &OrderBy{
			Direction: AscOrder,
			Field:     field,
		}
	}

	return &OrderBy{
		Direction: direction,
		Field:     field,
	}
}

func ParseOrderBy(allowedFields []string, orderby string, defaultOrder *OrderBy) (*OrderBy, error) {
	if orderby == "" {
		return defaultOrder, nil
	}

	var direction string

	if strings.HasPrefix(orderby, "-") {
		direction = DescOrder
		orderby = orderby[1:]
	} else {
		direction = AscOrder
	}

	if !slices.Contains(allowedFields, orderby) {
		return nil, fmt.Errorf("field %s does not exist", orderby)
	}

	return NewOrderBy(orderby, direction), nil
}
