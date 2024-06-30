package outletrepo

import (
	"time"

	"github.com/google/uuid"
	"github.com/goplateframework/internal/domain/address"
	"github.com/goplateframework/internal/domain/outlet"
)

type Model struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Phone       string    `db:"phone"`
	OpeningTime string    `db:"opening_time"`
	ClosingTime string    `db:"closing_time"`
	AddressID   uuid.UUID `db:"address_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func intoModel(o *outlet.OutletDTO) *Model {
	return &Model{
		ID:          o.ID,
		Name:        o.Name,
		Phone:       o.Phone,
		OpeningTime: o.OpeningTime,
		ClosingTime: o.ClosingTime,
		AddressID:   o.Address.ID,
		UpdatedAt:   o.UpdatedAt,
		CreatedAt:   o.CreatedAt,
	}
}

type ModelWithAddress struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Phone       string    `db:"phone"`
	OpeningTime string    `db:"opening_time"`
	ClosingTime string    `db:"closing_time"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	AddressID   uuid.UUID `db:"address_id"`

	Street     string `db:"street"`
	City       string `db:"city"`
	Province   string `db:"province"`
	PostalCode string `db:"postal_code"`
}

func (ma *ModelWithAddress) intoDTO() *outlet.OutletDTO {
	return &outlet.OutletDTO{
		ID:          ma.ID,
		Name:        ma.Name,
		Phone:       ma.Phone,
		OpeningTime: ma.OpeningTime,
		ClosingTime: ma.ClosingTime,
		Address: &address.AddressDTO{
			ID:         ma.AddressID,
			Street:     ma.Street,
			City:       ma.City,
			Province:   ma.Province,
			PostalCode: ma.PostalCode,
		},
	}
}
