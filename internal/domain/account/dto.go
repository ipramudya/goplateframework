package account

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/goplateframework/internal/sdk/validate"
)

type AccountDTO struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Password  string `json:"-"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
}

type AccountWithTokenDTO struct {
	Token   string      `json:"token"`
	Account *AccountDTO `json:"account"`
}

type NewAccouuntDTO struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

func (na NewAccouuntDTO) Validate() error {
	return validation.ValidateStruct(
		&na,
		validation.Field(&na.Firstname, validation.Required, validation.Length(1, 30)),
		validation.Field(&na.Lastname, validation.Required, validation.Length(1, 30)),
		validation.Field(&na.Email, validation.Required, is.Email),
		validation.Field(&na.Password, validation.Required),
		validation.Field(&na.Phone, validation.Required, validate.Phone),
	)
}

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

type ChangePasswordDTO struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (d ChangePasswordDTO) Validate() error {
	return validation.ValidateStruct(
		&d,
		validation.Field(&d.OldPassword, validation.Required),
		validation.Field(&d.NewPassword, validation.Required),
	)
}
