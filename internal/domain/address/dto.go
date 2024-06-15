package address

import validation "github.com/go-ozzo/ozzo-validation/v4"

type AddressDTO struct {
	ID         string `json:"id"`
	Street     string `json:"street"`
	City       string `json:"city"`
	Province   string `json:"province"`
	PostalCode string `json:"postal_code"`
}

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
