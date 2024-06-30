package addressx

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

// AddressDTO is what we send to client
type AddressDTO struct {
	ID         uuid.UUID `json:"id"`
	Street     string    `json:"street"`
	City       string    `json:"city"`
	Province   string    `json:"province"`
	PostalCode string    `json:"postal_code"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// NewAddressDTO is what client should send to create new address
type NewAddressDTO struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	Province   string `json:"province"`
	PostalCode string `json:"postal_code"`
}

func (n NewAddressDTO) Validate() error {
	return validation.ValidateStruct(
		&n,
		validation.Field(&n.Street, validation.Required, validation.Length(1, 255)),
		validation.Field(&n.City, validation.Required, validation.Length(1, 50)),
		validation.Field(&n.Province, validation.Required, validation.Length(1, 50)),
		validation.Field(&n.PostalCode, validation.Required, validation.Length(1, 10)),
	)
}
