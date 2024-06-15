package outlet

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/goplateframework/internal/domain/address"
	"github.com/goplateframework/internal/sdk/validate"
)

type OutletDTO struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Phone       string              `json:"phone"`
	OpeningTime string              `json:"opening_time"`
	ClosingTime string              `json:"closing_time"`
	Address     *address.AddressDTO `json:"address"`
}

type NewOutletDTO struct {
	Name        string                 `json:"name"`
	Phone       string                 `json:"phone"`
	OpeningTime string                 `json:"opening_time"`
	ClosingTime string                 `json:"closing_time"`
	Address     *address.NewAddressDTO `json:"address"`
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
