package menuxuc

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/domain/menux"
	"github.com/goplateframework/internal/domain/menux/menuxweb"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/pkg/logger"
)

// required repository methods which this usecase needs to store or retrieve data
type Repository interface {
	Create(ctx context.Context, nm *menux.MenuDTO) error
	GetAll(ctx context.Context, qp *menuxweb.QueryParams) (*[]menux.MenuDTO, error)
	Update(ctx context.Context, nm *menux.MenuDTO) error
}

type Usecase struct {
	conf       *config.Config
	log        *logger.Log
	menuDBRepo Repository
}

func New(conf *config.Config, log *logger.Log, menuDBRepo Repository) *Usecase {
	return &Usecase{
		conf:       conf,
		log:        log,
		menuDBRepo: menuDBRepo,
	}
}

func (uc *Usecase) Create(ctx context.Context, nm *menux.NewMenuDTO) (*menux.MenuDTO, error) {
	now := time.Now()

	m := &menux.MenuDTO{
		ID:          uuid.New().String(),
		Name:        nm.Name,
		Description: nm.Description,
		Price:       nm.Price,
		IsAvailable: nm.IsAvailable,
		ImageURL:    "",
		OutletID:    nm.OutletID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := uc.menuDBRepo.Create(ctx, m)

	if err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Errorf(e.DebugWithDetail(err.Error()))
		return nil, e
	}

	return m, nil
}

func (uc *Usecase) GetAll(ctx context.Context, qp *menuxweb.QueryParams) (*[]menux.MenuDTO, error) {
	m, err := uc.menuDBRepo.GetAll(ctx, qp)

	if err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.DebugWithDetail(err.Error()))
		return nil, e
	}

	return m, nil
}

func (uc *Usecase) Update(ctx context.Context, nm *menux.NewMenuDTO, id string) (*menux.MenuDTO, error) {
	m := &menux.MenuDTO{
		ID:          id,
		Name:        nm.Name,
		Description: nm.Description,
		Price:       nm.Price,
		IsAvailable: nm.IsAvailable,
		ImageURL:    "",
		OutletID:    nm.OutletID,
		UpdatedAt:   time.Now(),
	}

	err := uc.menuDBRepo.Update(ctx, m)

	if err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.DebugWithDetail(err.Error()))
		return nil, e
	}

	return m, nil
}
