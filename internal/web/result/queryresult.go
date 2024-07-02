package result

import "math"

type Result[T any] struct {
	Items   []T `json:"items"`
	Total   int `json:"total"`
	Page    int `json:"page"`
	MaxPage int `json:"max_page"`
	Size    int `json:"size_per_page"`
}

func New[T any](items []T, total, page, size int) *Result[T] {
	if items == nil {
		items = []T{}
	}

	return &Result[T]{
		Items:   items,
		Total:   total,
		Page:    page,
		MaxPage: int(math.Ceil(float64(total) / float64(size))),
		Size:    size,
	}
}
