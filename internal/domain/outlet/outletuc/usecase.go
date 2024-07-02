package outletuc

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/domain/address"
	"github.com/goplateframework/internal/domain/outlet"
	"github.com/goplateframework/internal/domain/outlet/outletweb"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/internal/web/result"
	"github.com/goplateframework/pkg/logger"
)

type iRepository interface {
	GetAll(ctx context.Context, qp *outletweb.QueryParams) ([]outlet.OutletDTO, error)
	GetOne(ctx context.Context, id uuid.UUID) (*outlet.OutletDTO, error)
	Create(ctx context.Context, a *outlet.OutletDTO) error
	Update(ctx context.Context, o *outlet.OutletDTO) error
	Count(ctx context.Context) (int, error)
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

func (uc *Usecase) Create(ctx context.Context, no *outlet.NewOutletDTO) (*outlet.OutletDTO, error) {
	now := time.Now()

	o := &outlet.OutletDTO{
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

func (uc *Usecase) GetAll(ctx context.Context, qp *outletweb.QueryParams) (*result.Result[outlet.OutletDTO], error) {
	total, err := uc.repo.Count(ctx)
	if err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.DebugWithDetail(err.Error()))
		return nil, e
	}

	if !qp.Page.CanPaginate(total) {
		e := errs.New(errs.InvalidArgument, errors.New("page requested is out of range"))
		uc.log.Error(e.Debug())
		return nil, e
	}

	o, err := uc.repo.GetAll(ctx, qp)
	if err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.DebugWithDetail(err.Error()))
		return nil, e
	}

	return result.New(o, total, qp.Page.Number, qp.Page.Size), nil
}

func (uc *Usecase) GetOne(ctx context.Context, id uuid.UUID) (*outlet.OutletDTO, error) {
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

func (uc *Usecase) Update(ctx context.Context, no *outlet.NewOutletDTO, id uuid.UUID) (*outlet.OutletDTO, error) {
	o := &outlet.OutletDTO{
		ID:          id,
		Name:        no.Name,
		Phone:       no.Phone,
		OpeningTime: no.OpeningTime,
		ClosingTime: no.ClosingTime,
		UpdatedAt:   time.Now(),
		Address:     &address.AddressDTO{},
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
