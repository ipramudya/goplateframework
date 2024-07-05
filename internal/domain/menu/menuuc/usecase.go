package menuuc

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/domain/menu"
	"github.com/goplateframework/internal/domain/menu/menuweb"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/internal/web/result"
	"github.com/goplateframework/internal/worker/pb"
	"github.com/goplateframework/pkg/logger"
)

// required iRepository methods which this usecase needs to store or retrieve data
type iRepository interface {
	Create(ctx context.Context, nm *menu.MenuDTO) error
	GetAll(ctx context.Context, qp *menuweb.QueryParams) ([]menu.MenuDTO, error)
	Update(ctx context.Context, nm *menu.MenuDTO) error
	Delete(ctx context.Context, id uuid.UUID) error
	Count(ctx context.Context, outletId string) (int, error)
}

type Usecase struct {
	conf       *config.Config
	log        *logger.Log
	menuDBRepo iRepository
	worker     pb.WorkerClient
}

func New(conf *config.Config, log *logger.Log, worker pb.WorkerClient, menuDBRepo iRepository) *Usecase {
	return &Usecase{
		conf:       conf,
		log:        log,
		menuDBRepo: menuDBRepo,
		worker:     worker,
	}
}

func (uc *Usecase) Create(ctx context.Context, nm *menu.NewMenuDTO, menuImage []byte) (*menu.MenuDTO, error) {
	now := time.Now()

	id := uuid.New()
	m := &menu.MenuDTO{
		ID:          id,
		Name:        nm.Name,
		Description: nm.Description,
		Price:       nm.Price,
		IsAvailable: nm.IsAvailable,
		ImageURL:    "pending",
		OutletID:    nm.OutletID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := uc.menuDBRepo.Create(ctx, m); err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Errorf(e.DebugWithDetail(err.Error()))
		return nil, e
	}

	go func() {
		workerCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := uc.worker.ProcessImage(workerCtx, &pb.ProcessImageRequest{
			Id:        id.String(),
			Table:     "menus",
			ImageData: menuImage,
		})

		if err != nil {
			uc.log.Error(err.Error())
		}
	}()

	return m, nil
}

func (uc *Usecase) GetAll(ctx context.Context, qp *menuweb.QueryParams) (*result.Result[menu.MenuDTO], error) {
	total, err := uc.menuDBRepo.Count(ctx, qp.Filter.OutletId)
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

	m, err := uc.menuDBRepo.GetAll(ctx, qp)
	if err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.DebugWithDetail(err.Error()))
		return nil, e
	}

	return result.New(m, total, qp.Page.Number, qp.Page.Size), nil
}

func (uc *Usecase) Update(ctx context.Context, nm *menu.NewMenuDTO, id uuid.UUID) (*menu.MenuDTO, error) {
	m := &menu.MenuDTO{
		ID:          id,
		Name:        nm.Name,
		Description: nm.Description,
		Price:       nm.Price,
		IsAvailable: nm.IsAvailable,
		ImageURL:    "",
		OutletID:    nm.OutletID,
		UpdatedAt:   time.Now(),
	}

	if err := uc.menuDBRepo.Update(ctx, m); err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.DebugWithDetail(err.Error()))
		return nil, e
	}

	return m, nil
}

func (uc *Usecase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := uc.menuDBRepo.Delete(ctx, id); err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.DebugWithDetail(err.Error()))
		return e
	}

	return nil
}
