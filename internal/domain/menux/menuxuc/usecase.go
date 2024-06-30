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
	"github.com/goplateframework/internal/web/queryparams"
	"github.com/goplateframework/internal/web/result"
	"github.com/goplateframework/pkg/logger"
)

// required iRepository methods which this usecase needs to store or retrieve data
type iRepository interface {
	Create(ctx context.Context, nm *menux.MenuDTO) error
	GetAll(ctx context.Context, qp *menuxweb.QueryParams) ([]menux.MenuDTO, error)
	Update(ctx context.Context, nm *menux.MenuDTO) error
	Count(ctx context.Context, outletId string) (int, error)
}

type Usecase struct {
	conf       *config.Config
	log        *logger.Log
	menuDBRepo iRepository
}

func New(conf *config.Config, log *logger.Log, menuDBRepo iRepository) *Usecase {
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

func (uc *Usecase) GetAll(ctx context.Context, qp *menuxweb.QueryParams) (*result.Result[menux.MenuDTO], error) {
	total, err := uc.menuDBRepo.Count(ctx, qp.Filter.OutletId)
	if err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.DebugWithDetail(err.Error()))
		return nil, e
	}

	if !queryparams.IsAllowedPaging(total, qp.Page) {
		e := errs.New(errs.InvalidArgument, errors.New("page requested is out of range"))
		uc.log.Error(e.Debug())
		return nil, e
	}

	m, err := uc.menuDBRepo.GetAll(ctx, qp)
	if err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.DebugWithDetail(err.Error()))
		return nil, e
	}

	lastId := ""
	if len(m) > 0 {
		lastId = m[len(m)-1].ID
	}

	return result.New(m, total, qp.Page.Number, qp.Page.Size, lastId), nil
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
