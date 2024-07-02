package outletrepo

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/goplateframework/internal/domain/outlet/outletweb"
)

func (dbrepo *repository) useFilter(args map[string]any, qp *outletweb.QueryParams) []byte {
	filterBuf := bytes.NewBufferString("")
	var writableStrings []string

	if qp.Filter.Name != "" {
		args["name"] = fmt.Sprintf("%%%s%%", qp.Filter.Name)
		writableStrings = append(writableStrings, "name LIKE :name")
	}

	if qp.Filter.Operate != "" {
		// we dont need to create properties on args, just because we dont append placeholder into query
		// it can be done directly by appending query string

		if qp.Filter.Operate == "open" {
			writableStrings = append(writableStrings,
				" current_timestamp AT TIME ZONE 'Asia/Jakarta' >= opening_time ",
			)
		} else if qp.Filter.Operate == "close" {
			writableStrings = append(writableStrings,
				" current_timestamp AT TIME ZONE 'Asia/Jakarta' < opening_time ",
			)
		}
	}

	if len(writableStrings) > 0 {
		filterBuf.WriteString(" WHERE ")
		filterBuf.WriteString(strings.Join(writableStrings, " AND "))
	}

	return filterBuf.Bytes()
}
