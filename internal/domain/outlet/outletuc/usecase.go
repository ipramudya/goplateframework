package outletuc

import (
	"context"

	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/domain/outlet"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/pkg/logger"
)

type Usecase struct {
	conf   *config.Config
	log    *logger.Log
	repoDB outlet.DBRepository
}

func New(conf *config.Config, log *logger.Log, repo outlet.DBRepository) *Usecase {
	return &Usecase{
		conf:   conf,
		log:    log,
		repoDB: repo,
	}
}

func (uc *Usecase) Create(ctx context.Context, no *outlet.NewOutletDTO) (*outlet.OutletDTO, error) {
	oa, err := uc.repoDB.AddOne(ctx, no)

	if err != nil {
		e := errs.Newf(errs.Internal, "something went wrong!")
		uc.log.Error(e.Debug())
		return &outlet.OutletDTO{}, e
	}

	return oa.IntoOutletDTO(), nil
}

func (uc *Usecase) GetOne(ctx context.Context, id string) (*outlet.OutletDTO, error) {
	o, err := uc.repoDB.GetOneByID(ctx, id)

	if err != nil {
		e := errs.Newf(errs.Internal, "something went wrong!")
		uc.log.Error(e.Debug())
		return &outlet.OutletDTO{}, e
	}

	return o.IntoOutletDTO(), nil
}
