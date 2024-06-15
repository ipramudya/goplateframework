package addressuc

import (
	"context"
	"database/sql"
	"errors"

	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/domain/address"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/pkg/logger"
)

type Usecase struct {
	conf   *config.Config
	log    *logger.Log
	dbRepo address.DBRepository
}

func New(conf *config.Config, log *logger.Log, dbRepo address.DBRepository) *Usecase {
	return &Usecase{
		conf:   conf,
		log:    log,
		dbRepo: dbRepo,
	}
}

func (uc *Usecase) Update(ctx context.Context, na *address.NewAddressDTO, id string) (*address.AddressDTO, error) {
	a, err := uc.dbRepo.Update(ctx, na, id)

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

	return a.IntoAddressDTO(), nil
}
