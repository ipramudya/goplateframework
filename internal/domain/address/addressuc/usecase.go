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

func (uc *Usecase) Update(ctx context.Context, a *address.AddressDTO) error {
	return uc.repoDB.Update(ctx, a)
}
