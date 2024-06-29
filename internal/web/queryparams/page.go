package queryparams

import (
	"fmt"
	"math"
	"strconv"
)

type Page struct {
	Number int
	Size   int
}

const (
	DefaultNumber = 1
	DefaultSize   = 10
)

func ParsePage(pageNumber, pageSize string) (*Page, error) {
	number := DefaultNumber

	if pageNumber != "" {
		if n, err := strconv.Atoi(pageNumber); err != nil {
			return nil, err
		} else {
			number = n
		}
	}

	size := DefaultSize
	if pageSize != "" {
		if s, err := strconv.Atoi(pageSize); err != nil {
			return nil, err
		} else {
			size = s
		}
	}

	if number <= 0 {
		return nil, fmt.Errorf("page number must be greater than zero")
	}

	if size <= 0 {
		return nil, fmt.Errorf("page size must be greater than zero")
	}

	if size > 100 {
		return nil, fmt.Errorf("page size too big, max is 100")
	}

	return &Page{Number: number, Size: size}, nil
}

func IsAllowedPaging(total int, page *Page) bool {
	maxPage := int(math.Ceil(float64(total) / float64(page.Size)))
	return page.Number <= maxPage
}
