package address

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type AddressDTO struct {
	ID         string `json:"id"`
	Street     string `json:"street"`
	City       string `json:"city"`
	Province   string `json:"province"`
	PostalCode string `json:"postal_code"`
}

func (a *AddressDTO) Validate() error {
	return validation.ValidateStruct(
		&a,
		validation.Field(&a.Street, validation.Required, validation.Length(1, 255)),
		validation.Field(&a.City, validation.Required, validation.Length(1, 50)),
		validation.Field(&a.Province, validation.Required, validation.Length(1, 50)),
		validation.Field(&a.PostalCode, validation.Required, validation.Length(1, 10)),
	)
}

type NewAddressDTO struct {
	AddressDTO
}

func (a *NewAddressDTO) Validate() error {
	return validation.ValidateStruct(
		&a,
		validation.Field(&a.ID, validation.Required, validation.Length(1, 36), is.UUID),
		validation.Field(&a.Street, validation.Required, validation.Length(1, 255)),
		validation.Field(&a.City, validation.Required, validation.Length(1, 50)),
		validation.Field(&a.Province, validation.Required, validation.Length(1, 50)),
		validation.Field(&a.PostalCode, validation.Required, validation.Length(1, 10)),
	)
}
