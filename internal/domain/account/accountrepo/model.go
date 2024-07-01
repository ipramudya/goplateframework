package accountrepo

import (
	"time"

	"github.com/google/uuid"
	"github.com/goplateframework/internal/domain/account"
)

type Model struct {
	ID        uuid.UUID `db:"id"`
	Firstname string    `db:"firstname"`
	Lastname  string    `db:"lastname"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Phone     string    `db:"phone"`
	Role      string    `db:"role"` // ENUM: user, admin, superadmin
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (m *Model) intoDTO() *account.AccountDTO {
	return &account.AccountDTO{
		ID:        m.ID,
		Firstname: m.Firstname,
		Lastname:  m.Lastname,
		Password:  m.Password,
		Email:     m.Email,
		Phone:     m.Phone,
		Role:      m.Role,
	}
}

func intoModel(a *account.AccountDTO) *Model {
	return &Model{
		ID:        a.ID,
		Firstname: a.Firstname,
		Lastname:  a.Lastname,
		Email:     a.Email,
		Password:  a.Password,
		Phone:     a.Phone,
		Role:      a.Role,
	}
}
