package addressrepo

import (
	"time"

	"github.com/google/uuid"
	"github.com/goplateframework/internal/domain/address"
)

type Model struct {
	ID         uuid.UUID `db:"id"`
	Street     string    `db:"street"`
	City       string    `db:"city"`
	Province   string    `db:"province"`
	PostalCode string    `db:"postal_code"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func intoModel(a *address.AddressDTO) *Model {
	return &Model{
		ID:         a.ID,
		Street:     a.Street,
		City:       a.City,
		Province:   a.Province,
		PostalCode: a.PostalCode,
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
	}
}
