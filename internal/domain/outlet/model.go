package outlet

import (
	"time"

	"github.com/goplateframework/internal/domain/address"
)

type Schema struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Phone       string    `db:"phone"`
	OpeningTime string    `db:"opening_time"`
	ClosingTime string    `db:"closing_time"`
	AddressID   string    `db:"address_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type SchemaWithAddress struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Phone       string    `db:"phone"`
	OpeningTime string    `db:"opening_time"`
	ClosingTime string    `db:"closing_time"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	AddressID   string    `db:"address_id"`

	Street     string `db:"street"`
	City       string `db:"city"`
	Province   string `db:"province"`
	PostalCode string `db:"postal_code"`
}

func (s *SchemaWithAddress) IntoOutletDTO() *OutletDTO {
	return &OutletDTO{
		ID:          s.ID,
		Name:        s.Name,
		Phone:       s.Phone,
		OpeningTime: s.OpeningTime,
		ClosingTime: s.ClosingTime,
		Address: &address.AddressDTO{
			ID:         s.AddressID,
			Street:     s.Street,
			City:       s.City,
			Province:   s.Province,
			PostalCode: s.PostalCode,
		},
	}
}
