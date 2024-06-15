package outletuc

import (
	"context"
	"database/sql"
	"errors"

	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/domain/address"
	"github.com/goplateframework/internal/domain/outlet"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/pkg/logger"
)

type Usecase struct {
	conf         *config.Config
	log          *logger.Log
	outletDBRepo outlet.DBRepository
	addrDBRepo   address.DBRepository
}

func New(conf *config.Config, log *logger.Log, outletDBRepo outlet.DBRepository, addrDBRepo address.DBRepository) *Usecase {
	return &Usecase{
		conf:         conf,
		log:          log,
		outletDBRepo: outletDBRepo,
		addrDBRepo:   addrDBRepo,
	}
}

func (uc *Usecase) Create(ctx context.Context, no *outlet.NewOutletDTO) (*outlet.OutletDTO, error) {
	oa, err := uc.outletDBRepo.AddOne(ctx, no)

	if err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.Debug())
		return nil, e
	}

	return oa.IntoOutletDTO(), nil
}

func (uc *Usecase) GetOne(ctx context.Context, id string) (*outlet.OutletDTO, error) {
	o, err := uc.outletDBRepo.GetOneByID(ctx, id)

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

	return o.IntoOutletDTO(), nil
}

func (uc *Usecase) Update(ctx context.Context, no *outlet.NewOutletDTO, id string) (*outlet.OutletDTO, error) {
	o, err := uc.outletDBRepo.Update(ctx, no, id)

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

	a, err := uc.addrDBRepo.GetOneByID(ctx, o.AddressID)
	if err != nil {
		if err == sql.ErrNoRows {
			e := errs.New(errs.NotFound, errors.New("address not found"))
			uc.log.Error(e.Debug())
			return nil, e
		}

		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.Debug())
		return nil, e
	}

	return &outlet.OutletDTO{
		ID:          o.ID,
		Name:        o.Name,
		Phone:       o.Phone,
		OpeningTime: o.OpeningTime,
		ClosingTime: o.ClosingTime,
		Address: &address.AddressDTO{
			ID:         o.AddressID,
			Street:     a.Street,
			City:       a.City,
			Province:   a.Province,
			PostalCode: a.PostalCode,
		},
	}, nil
}
