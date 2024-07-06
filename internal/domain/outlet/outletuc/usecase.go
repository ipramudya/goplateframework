package outletuc

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/domain/address"
	"github.com/goplateframework/internal/domain/outlet"
	"github.com/goplateframework/internal/domain/outlet/outletweb"
	"github.com/goplateframework/internal/sdk/errshttp"
	"github.com/goplateframework/internal/web/result"
	"github.com/goplateframework/pkg/logger"
)

type iRepository interface {
	GetAll(ctx context.Context, qp *outletweb.QueryParams) ([]outlet.OutletDTO, error)
	GetOne(ctx context.Context, id uuid.UUID) (*outlet.OutletDTO, error)
	Create(ctx context.Context, a *outlet.OutletDTO) error
	Update(ctx context.Context, o *outlet.OutletDTO) error
	Delete(ctx context.Context, id uuid.UUID) error
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
		return nil, errshttp.New(errshttp.Internal, "Something went wrong")
	}

	return o, nil
}

func (uc *Usecase) GetAll(ctx context.Context, qp *outletweb.QueryParams) (*result.Result[outlet.OutletDTO], error) {
	total, err := uc.repo.Count(ctx)
	if err != nil {
		return nil, errshttp.New(errshttp.Internal, "Something went wrong")
	}

	if !qp.Page.CanPaginate(total) {
		e := errshttp.New(errshttp.InvalidArgument, "Page requested is out of range")
		e.AddDetail(fmt.Sprintf("pagination: page number must be between 1 and %d", total))
		return nil, e
	}

	o, err := uc.repo.GetAll(ctx, qp)
	if err != nil {
		return nil, errshttp.New(errshttp.Internal, "Something went wrong")
	}

	return result.New(o, total, qp.Page.Number, qp.Page.Size), nil
}

func (uc *Usecase) GetOne(ctx context.Context, id uuid.UUID) (*outlet.OutletDTO, error) {
	o, err := uc.repo.GetOne(ctx, id)

	if err != nil {
		if err == sql.ErrNoRows {
			e := errshttp.New(errshttp.NotFound, "Outlet not found")
			e.AddDetail(fmt.Sprintf("data: outlet with id %s not found", id))
			return nil, e
		}

		return nil, errshttp.New(errshttp.Internal, "Something went wrong")
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
		return nil, errshttp.New(errshttp.Internal, "Something went wrong")
	}

	oa, err := uc.repo.GetOne(ctx, id)
	if err != nil {
		return nil, errshttp.New(errshttp.Internal, "Something went wrong")
	}

	return oa, nil
}

func (uc *Usecase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return errshttp.New(errshttp.Internal, "Something went wrong")
	}

	return nil
}
