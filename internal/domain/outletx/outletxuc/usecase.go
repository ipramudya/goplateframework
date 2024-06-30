package outletxuc

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/domain/addressx"
	"github.com/goplateframework/internal/domain/outletx"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/pkg/logger"
)

type iRepository interface {
	GetOne(ctx context.Context, id uuid.UUID) (*outletx.OutletDTO, error)
	Create(ctx context.Context, a *outletx.OutletDTO) error
	Update(ctx context.Context, o *outletx.OutletDTO) error
}

type Usecase struct {
	conf *config.Config
	log  *logger.Log
	repo iRepository
}

func New(conf *config.Config, log *logger.Log, repo iRepository) *Usecase {
	return &Usecase{
		conf: conf,
		log:  log,
		repo: repo,
	}
}

func (uc *Usecase) Create(ctx context.Context, no *outletx.NewOutletDTO) (*outletx.OutletDTO, error) {
	now := time.Now()

	o := &outletx.OutletDTO{
		ID:          uuid.New(),
		Name:        no.Name,
		Phone:       no.Phone,
		OpeningTime: no.OpeningTime,
		ClosingTime: no.ClosingTime,
		CreatedAt:   now,
		UpdatedAt:   now,
		Address:     no.Address,
	}

	if err := uc.repo.Create(ctx, o); err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.Debug())
		return nil, e
	}

	return o, nil
}

func (uc *Usecase) GetOne(ctx context.Context, id uuid.UUID) (*outletx.OutletDTO, error) {
	o, err := uc.repo.GetOne(ctx, id)

	if err != nil {
		if err == sql.ErrNoRows {
			e := errs.New(errs.NotFound, errors.New("outlet not found"))
			uc.log.Error(e.Debug())
			return nil, e
		}

		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.Debug())
		return nil, e
	}

	return o, nil
}

func (uc *Usecase) Update(ctx context.Context, no *outletx.NewOutletDTO, id uuid.UUID) (*outletx.OutletDTO, error) {
	o := &outletx.OutletDTO{
		ID:          id,
		Name:        no.Name,
		Phone:       no.Phone,
		OpeningTime: no.OpeningTime,
		ClosingTime: no.ClosingTime,
		UpdatedAt:   time.Now(),
		Address:     &addressx.AddressDTO{},
	}

	if err := uc.repo.Update(ctx, o); err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.Debug())
		return nil, e
	}

	oa, err := uc.repo.GetOne(ctx, id)
	if err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.Debug())
		return nil, e
	}

	return oa, nil
}
