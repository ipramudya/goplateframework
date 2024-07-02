package outletrepo

import (
	"fmt"
	"strings"

	"github.com/goplateframework/internal/domain/outlet/outletweb"
)

func (dbrepo *repository) buildFilter(args map[string]any, qp *outletweb.QueryParams) string {
	var filters []string

	if qp.Filter.Name != "" {
		args["name"] = "%" + qp.Filter.Name + "%"
		filters = append(filters, " name ILIKE :name")
	}

	if qp.Filter.Operate != "" {
		// we dont need to create properties on args, just because we dont append placeholder into query
		// it can be done directly by appending query string

		if qp.Filter.Operate == "open" {
			filters = append(filters, " current_timestamp AT TIME ZONE 'Asia/Jakarta' > opening_time")
		} else if qp.Filter.Operate == "close" {
			filters = append(filters, " current_timestamp AT TIME ZONE 'Asia/Jakarta' < opening_time")
		}
	}

	if len(filters) > 0 {
		return fmt.Sprintf(" WHERE %s", strings.Join(filters, " AND "))
	}

	return ""
}
