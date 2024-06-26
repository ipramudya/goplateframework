package menuuc

import (
	"context"
	"database/sql"
	"errors"

	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/domain/menu"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/pkg/logger"
)

type Usecase struct {
	conf       *config.Config
	log        *logger.Log
	menuDBRepo menu.DBRepository
}

func New(conf *config.Config, log *logger.Log, menuDBRepo menu.DBRepository) *Usecase {
	return &Usecase{
		conf:       conf,
		log:        log,
		menuDBRepo: menuDBRepo,
	}
}

func (uc *Usecase) Create(ctx context.Context, nm *menu.NewMenuDTO) (*menu.MenuDTO, error) {
	m, err := uc.menuDBRepo.AddOne(ctx, nm)

	if err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.Debug())
		return nil, e
	}

	return m.IntoMenuDTO(), nil
}

func (uc *Usecase) Update(ctx context.Context, nm *menu.NewMenuDTO, id string) (*menu.MenuDTO, error) {
	m, err := uc.menuDBRepo.Update(ctx, nm, id)

	if err != nil {
		if err == sql.ErrNoRows {
			e := errs.New(errs.NotFound, errors.New("menu not found"))
			uc.log.Error(e.Debug())
			return nil, e
		}

		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.Debug())
		return nil, e
	}

	return m.IntoMenuDTO(), nil
}

func (uc *Usecase) GetAllByOutletID(ctx context.Context, outletID string) ([]*menu.MenuDTO, error) {
	m, err := uc.menuDBRepo.GetAllByOutletID(ctx, outletID)

	if err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.Debug())
		return nil, e
	}

	menus := make([]*menu.MenuDTO, 0)
	for _, v := range *m {
		menus = append(menus, v.IntoMenuDTO())
	}

	return menus, nil
}
