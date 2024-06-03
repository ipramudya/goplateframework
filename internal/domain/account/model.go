package account

import (
	"time"

	"github.com/google/uuid"
)

type Schema struct {
	ID        uuid.UUID `db:"id"`
	Firstname string    `db:"firstname"`
	Lastname  string    `db:"lastname"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Phone     string    `db:"phone"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (s *Schema) IntoAccountDTO() *AccountDTO {
	return &AccountDTO{
		ID:        s.ID.String(),
		Firstname: s.Firstname,
		Lastname:  s.Lastname,
		Password:  s.Password,
		Email:     s.Email,
		Phone:     s.Phone,
		Role:      s.Role,
	}
}

func (s *Schema) IntoAccountWithTokenDTO(token string) *AccountWithTokenDTO {
	return &AccountWithTokenDTO{
		Account: s.IntoAccountDTO(),
		Token:   token,
	}
}
