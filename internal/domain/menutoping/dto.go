package menutoping

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

type MenuTopingsDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	IsAvailable bool      `json:"is_available"`
	ImageURL    string    `json:"image_url"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	MenuID      uuid.UUID `json:"menu_id"`
}

type NewMenuTopingsDTO struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	IsAvailable bool    `json:"is_available"`
	Stock       int     `json:"stock"`
	MenuID      uuid.UUID
	// ImageURL    string  `json:"image_url"`
}

func (nmt NewMenuTopingsDTO) Validate() error {
	return validation.ValidateStruct(&nmt,
		validation.Field(&nmt.Name, validation.Required, validation.Length(1, 50)),
		validation.Field(&nmt.Price, validation.Required),
		validation.Field(&nmt.IsAvailable, validation.Required),
		validation.Field(&nmt.Stock, validation.Required),
		validation.Field(&nmt.MenuID, validation.Required, is.UUIDv4),
	)
}
