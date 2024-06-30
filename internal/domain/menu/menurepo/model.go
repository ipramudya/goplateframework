package menurepo

import (
	"time"

	"github.com/goplateframework/internal/domain/menu"
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

func intoModel(m *menu.MenuDTO) *Model {
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

func (m *Model) intoDTO() *menu.MenuDTO {
	return &menu.MenuDTO{
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
