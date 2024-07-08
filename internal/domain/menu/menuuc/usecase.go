package menuuc

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/domain/menu"
	"github.com/goplateframework/internal/domain/menu/menuweb"
	"github.com/goplateframework/internal/sdk/errshttp"
	"github.com/goplateframework/internal/web/result"
	"github.com/goplateframework/internal/worker/pb"
	"github.com/goplateframework/pkg/logger"
)

// required iRepository methods which this usecase needs to store or retrieve data
type iRepository interface {
	Create(ctx context.Context, nm *menu.MenuDTO) error
	GetAll(ctx context.Context, qp *menuweb.QueryParams) ([]menu.MenuDTO, error)
	GetOne(ctx context.Context, id uuid.UUID) (*menu.MenuDTO, error)
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

func (uc *Usecase) Create(ctx context.Context, nm *menu.NewMenuDTO, image *[]byte) (*menu.MenuDTO, error) {
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
		return nil, errshttp.New(errshttp.Internal, "Something went wrong")
	}

	go func() {
		workerCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := uc.worker.ProcessImage(workerCtx, &pb.ProcessImageRequest{
			Id:        id.String(),
			Table:     "menus",
			ImageData: *image,
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
		return nil, errshttp.New(errshttp.Internal, "Something went wrong")
	}

	if !qp.Page.CanPaginate(total) {
		e := errshttp.New(errshttp.InvalidArgument, "Page requested is out of range")
		e.AddDetail(fmt.Sprintf("pagination: page number must be between 1 and %d", total))
		return nil, e
	}

	m, err := uc.menuDBRepo.GetAll(ctx, qp)
	if err != nil {
		return nil, errshttp.New(errshttp.Internal, "Something went wrong")
	}

	return result.New(m, total, qp.Page.Number, qp.Page.Size), nil
}

func (uc *Usecase) Update(ctx context.Context, nm *menu.NewMenuDTO, id uuid.UUID, image *[]byte) (*menu.MenuDTO, error) {
	if image != nil {
		nm.ImageURL = "pending"
	}

	m := &menu.MenuDTO{
		ID:          id,
		Name:        nm.Name,
		Description: nm.Description,
		Price:       nm.Price,
		IsAvailable: nm.IsAvailable,
		ImageURL:    nm.ImageURL,
		OutletID:    nm.OutletID,
		UpdatedAt:   time.Now(),
	}

	if err := uc.menuDBRepo.Update(ctx, m); err != nil {
		return nil, errshttp.New(errshttp.Internal, "Something went wrong")
	}

	if image != nil {
		go func() {
			workerCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err := uc.worker.ProcessImage(workerCtx, &pb.ProcessImageRequest{
				Id:        id.String(),
				Table:     "menus",
				ImageData: *image,
			})

			if err != nil {
				uc.log.Error(err.Error())
			}
		}()
	}

	return m, nil
}

func (uc *Usecase) Delete(ctx context.Context, id uuid.UUID) error {
	m, err := uc.menuDBRepo.GetOne(ctx, id)

	if err == sql.ErrNoRows {
		return errshttp.New(errshttp.NotFound, "Menu topping not found")
	}

	if err := uc.menuDBRepo.Delete(ctx, id); err != nil {
		return errshttp.New(errshttp.Internal, "Something went wrong")
	}

	go func() {
		workerCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err = uc.worker.DeleteImage(workerCtx, &pb.DeleteImageRequest{
			Table:    "menus",
			ImageUrl: m.ImageURL,
		})

		if err != nil {
			uc.log.Error(fmt.Errorf("failed to delete image: %w", err).Error())
		}
	}()

	return nil
}
