package menu

import "time"

type Schema struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       float64   `db:"price"`
	IsAvailable bool      `db:"is_available"`
	ImageURL    string    `db:"image_url"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	OutletID    string    `db:"outlet_id"`
}

func (s *Schema) IntoMenuDTO() *MenuDTO {
	return &MenuDTO{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Price:       s.Price,
		IsAvailable: s.IsAvailable,
		ImageURL:    s.ImageURL,
	}
}
