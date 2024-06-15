package address

import (
	"time"

	"github.com/google/uuid"
)

type Schema struct {
	ID         uuid.UUID `db:"id"`
	Street     string    `db:"street"`
	City       string    `db:"city"`
	Province   string    `db:"province"`
	PostalCode string    `db:"postal_code"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
