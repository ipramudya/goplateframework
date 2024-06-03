package auth

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (d LoginDTO) Validate() error {
	return validation.ValidateStruct(
		&d,
		validation.Field(&d.Email, validation.Required, is.Email),
		validation.Field(&d.Password, validation.Required),
	)
}
