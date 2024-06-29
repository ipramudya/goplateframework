package menuxrepo

import (
	"time"

	"github.com/goplateframework/internal/domain/menux"
)

type Model struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       float64   `db:"price"`
	IsAvailable bool      `db:"is_available"`
	ImageURL    string    `db:"image_url"`
	OutletID    string    `db:"outlet_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func intoModel(m *menux.MenuDTO) *Model {
	return &Model{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Price:       m.Price,
		IsAvailable: m.IsAvailable,
		ImageURL:    m.ImageURL,
		OutletID:    m.OutletID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (m *Model) intoDTO() *menux.MenuDTO {
	return &menux.MenuDTO{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Price:       m.Price,
		IsAvailable: m.IsAvailable,
		ImageURL:    m.ImageURL,
		OutletID:    m.OutletID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
