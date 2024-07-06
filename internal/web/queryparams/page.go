package queryparams

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

type Page struct {
	Number int
	Offset int
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
		return nil, errors.New("pagination: page number must be greater than zero")
	}

	if size <= 0 {
		return nil, fmt.Errorf("pagination: size must be greater than zero")
	}

	if size > 100 {
		return nil, fmt.Errorf("pagination: page size too big, max is 100")
	}

	return &Page{
		Number: number,
		Size:   size,
		Offset: (number - 1) * size,
	}, nil
}

func (page *Page) CanPaginate(total int) bool {
	maxPage := int(math.Ceil(float64(total) / float64(page.Size)))
	return page.Number <= maxPage
}
