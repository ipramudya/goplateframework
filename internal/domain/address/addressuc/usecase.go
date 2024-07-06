package addressuc

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/domain/address"
	"github.com/goplateframework/internal/sdk/errshttp"
	"github.com/goplateframework/pkg/logger"
)

type iRepository interface {
	Create(ctx context.Context, a *address.AddressDTO) error
	GetOne(ctx context.Context, id uuid.UUID) (*address.AddressDTO, error)
	Update(ctx context.Context, na *address.AddressDTO) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Usecase struct {
	conf *config.Config
	log  *logger.Log
	repo iRepository
}

func New(conf *config.Config, log *logger.Log, addressDBRepo iRepository) *Usecase {
	return &Usecase{
		conf: conf,
		log:  log,
		repo: addressDBRepo,
	}
}

func (uc *Usecase) Create(ctx context.Context, na *address.NewAddressDTO) (*address.AddressDTO, error) {
	a := &address.AddressDTO{
		ID:         uuid.New(),
		Street:     na.Street,
		City:       na.City,
		Province:   na.Province,
		PostalCode: na.PostalCode,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := uc.repo.Create(ctx, a); err != nil {
		return nil, errshttp.New(errshttp.Internal, "Something went wrong")
	}

	return a, nil
}

func (uc *Usecase) GetOne(ctx context.Context, id uuid.UUID) (*address.AddressDTO, error) {
	a, err := uc.repo.GetOne(ctx, id)

	if err != nil {
		if err == sql.ErrNoRows {
			e := errshttp.New(errshttp.NotFound, "Could not find address")
			e.AddDetail(fmt.Sprintf("data: address with id %s not found", id))
			return nil, e
		}

		return nil, errshttp.New(errshttp.Internal, "Something went wrong")
	}

	return a, nil
}

func (uc *Usecase) Update(ctx context.Context, na *address.NewAddressDTO, id uuid.UUID) (*address.AddressDTO, error) {
	a := &address.AddressDTO{
		ID:         id,
		Street:     na.Street,
		City:       na.City,
		Province:   na.Province,
		PostalCode: na.PostalCode,
		UpdatedAt:  time.Now(),
	}

	if err := uc.repo.Update(ctx, a); err != nil {
		if err == sql.ErrNoRows {
			e := errshttp.New(errshttp.NotFound, "Could not find address")
			e.AddDetail(fmt.Sprintf("data: address with id %s not found", id))
			return nil, e
		}

		return nil, errshttp.New(errshttp.Internal, "Something went wrong")
	}

	return a, nil
}

func (uc *Usecase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return errshttp.New(errshttp.Internal, "Something went wrong")
	}

	return nil
}
