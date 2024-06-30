package outletx

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/goplateframework/internal/domain/addressx"
	"github.com/goplateframework/internal/sdk/validate"
)

type OutletDTO struct {
	ID          uuid.UUID            `json:"id"`
	Name        string               `json:"name"`
	Phone       string               `json:"phone"`
	OpeningTime string               `json:"opening_time"`
	ClosingTime string               `json:"closing_time"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	Address     *addressx.AddressDTO `json:"address"`
}

type NewOutletDTO struct {
	Name        string               `json:"name"`
	Phone       string               `json:"phone"`
	OpeningTime string               `json:"opening_time"`
	ClosingTime string               `json:"closing_time"`
	Address     *addressx.AddressDTO `json:"address"`
}

func (o NewOutletDTO) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.Name, validation.Required, validation.Length(1, 50)),
		validation.Field(&o.Phone, validation.Required, validate.Phone),
		validation.Field(&o.OpeningTime, validation.Required, validate.Timestamp),
		validation.Field(&o.ClosingTime, validation.Required, validate.Timestamp),
		validation.Field(&o.Address),
	)
}
