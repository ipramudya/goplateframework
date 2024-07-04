package menutopingrepo

import (
	"time"

	"github.com/google/uuid"
	"github.com/goplateframework/internal/domain/menutoping"
)

type Model struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Price       float64   `db:"price"`
	IsAvailable bool      `db:"is_available"`
	ImageURL    string    `db:"image_url"`
	Stock       int       `db:"stock"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	MenuID      uuid.UUID `db:"menu_id"`
}

func intoModel(mt *menutoping.MenuTopingsDTO) *Model {
	return &Model{
		ID:          mt.ID,
		Name:        mt.Name,
		Price:       mt.Price,
		IsAvailable: mt.IsAvailable,
		ImageURL:    mt.ImageURL,
		Stock:       mt.Stock,
		CreatedAt:   mt.CreatedAt,
		UpdatedAt:   mt.UpdatedAt,
		MenuID:      mt.MenuID,
	}
}

func (m *Model) intoDTO() *menutoping.MenuTopingsDTO {
	return &menutoping.MenuTopingsDTO{
		ID:          m.ID,
		Name:        m.Name,
		Price:       m.Price,
		IsAvailable: m.IsAvailable,
		ImageURL:    m.ImageURL,
		Stock:       m.Stock,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		MenuID:      m.MenuID,
	}
}
