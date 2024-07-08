package menu

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

// MenuDTO is what we send to client
type MenuDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	IsAvailable bool      `json:"is_available"`
	ImageURL    string    `json:"image_url"`
	OutletID    string    `json:"outlet_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewMenuDTO is what client should send to create new menu
type NewMenuDTO struct {
	Name        string  `json:"name" form:"name"`
	Description string  `json:"description" form:"description"`
	Price       float64 `json:"price" form:"price"`
	ImageURL    string  `json:"image_url" form:"image_url"`
	IsAvailable bool    `json:"is_available" form:"is_available"`
	OutletID    string  `json:"outlet_id" form:"outlet_id"`
}

func (m NewMenuDTO) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(1, 50)),
		validation.Field(&m.Description, validation.Required, validation.Length(1, 255)),
		validation.Field(&m.Price, validation.Required),
		validation.Field(&m.OutletID, validation.Required, is.UUIDv4),
	)
}
