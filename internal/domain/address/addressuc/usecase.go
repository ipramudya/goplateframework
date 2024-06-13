package addressuc

import (
	"context"

	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/domain/address"
	"github.com/goplateframework/pkg/logger"
)

type Usecase struct {
	conf   *config.Config
	log    *logger.Log
	repoDB address.DBRepository
}

func New(conf *config.Config, log *logger.Log, repo address.DBRepository) *Usecase {
	return &Usecase{
		conf:   conf,
		log:    log,
		repoDB: repo,
	}
}

func (uc *Usecase) AddOne(ctx context.Context, a *address.NewAddressDTO) error {
	if err := a.Validate(); err != nil {
		return err
	}
	return uc.repoDB.AddOne(ctx, a)
}

func (uc *Usecase) GetOneByID(ctx context.Context, id string) (*address.Schema, error) {
	return uc.repoDB.GetOneByID(ctx, id)
}

func (uc *Usecase) Update(ctx context.Context, a *address.NewAddressDTO) error {
	if err := a.Validate(); err != nil {
		return err
	}
	return uc.repoDB.Update(ctx, a)
}
