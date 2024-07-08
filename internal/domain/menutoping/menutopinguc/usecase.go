package menutopinguc

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/domain/menutoping"
	"github.com/goplateframework/internal/sdk/errshttp"
	"github.com/goplateframework/internal/worker/pb"
	"github.com/goplateframework/pkg/logger"
)

// required iRepository methods which this usecase needs to store or retrieve data
type iRepository interface {
	Create(ctx context.Context, m *menutoping.MenuTopingsDTO) error
	GetAll(ctx context.Context) ([]*menutoping.MenuTopingsDTO, error)
	GetOne(ctx context.Context, id uuid.UUID) (*menutoping.MenuTopingsDTO, error)
	Update(ctx context.Context, m *menutoping.MenuTopingsDTO) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Usecase struct {
	conf             *config.Config
	log              *logger.Log
	menuTopingDBRepo iRepository
	worker           pb.WorkerClient
}

func New(conf *config.Config, log *logger.Log, menuTopingDBRepo iRepository, worker pb.WorkerClient) *Usecase {
	return &Usecase{
		conf:             conf,
		log:              log,
		menuTopingDBRepo: menuTopingDBRepo,
		worker:           worker,
	}
}

func (uc *Usecase) Create(ctx context.Context, nmt *menutoping.NewMenuTopingsDTO, image *[]byte) (*menutoping.MenuTopingsDTO, error) {
	now := time.Now()
	id := uuid.New()

	mt := &menutoping.MenuTopingsDTO{
		ID:          id,
		Name:        nmt.Name,
		Price:       nmt.Price,
		IsAvailable: nmt.IsAvailable,
		ImageURL:    "pending",
		Stock:       nmt.Stock,
		CreatedAt:   now,
		UpdatedAt:   now,
		MenuID:      nmt.MenuID,
	}

	if err := uc.menuTopingDBRepo.Create(ctx, mt); err != nil {
		return nil, errshttp.New(errshttp.Internal, "Something went wrong")
	}

	go func() {
		workerCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := uc.worker.ProcessImage(workerCtx, &pb.ProcessImageRequest{
			Id:        id.String(),
			Table:     "menu_topings",
			ImageData: *image,
		})

		if err != nil {
			uc.log.Error(err.Error())
		}
	}()

	return mt, nil
}

func (uc *Usecase) GetAll(ctx context.Context) ([]*menutoping.MenuTopingsDTO, error) {
	mt, err := uc.menuTopingDBRepo.GetAll(ctx)

	if err != nil {
		return nil, errshttp.New(errshttp.Internal, "Something went wrong")
	}

	return mt, nil
}

func (uc *Usecase) GetOne(ctx context.Context, id uuid.UUID) (*menutoping.MenuTopingsDTO, error) {
	mt, err := uc.menuTopingDBRepo.GetOne(ctx, id)

	if err != nil {
		return nil, errshttp.New(errshttp.Internal, "Something went wrong")
	}

	return mt, nil
}

func (uc *Usecase) Update(ctx context.Context, nmt *menutoping.NewMenuTopingsDTO, id uuid.UUID, image *[]byte) (*menutoping.MenuTopingsDTO, error) {
	mt := &menutoping.MenuTopingsDTO{
		ID:          id,
		Name:        nmt.Name,
		Price:       nmt.Price,
		IsAvailable: nmt.IsAvailable,
		Stock:       nmt.Stock,
		UpdatedAt:   time.Now(),
	}

	if err := uc.menuTopingDBRepo.Update(ctx, mt); err != nil {
		return nil, errshttp.New(errshttp.Internal, "Something went wrong")
	}

	go func() {
		workerCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := uc.worker.ProcessImage(workerCtx, &pb.ProcessImageRequest{
			Id:        id.String(),
			Table:     "menu_topings",
			ImageData: *image,
		})

		if err != nil {
			uc.log.Error(err.Error())
		}
	}()

	return mt, nil
}

func (uc *Usecase) Delete(ctx context.Context, id uuid.UUID) error {
	mt, err := uc.menuTopingDBRepo.GetOne(ctx, id)

	if err == sql.ErrNoRows {
		return errshttp.New(errshttp.NotFound, "Menu topping not found")
	}

	if err := uc.menuTopingDBRepo.Delete(ctx, id); err != nil {
		return errshttp.New(errshttp.Internal, "Something went wrong")
	}

	go func() {
		workerCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err = uc.worker.DeleteImage(workerCtx, &pb.DeleteImageRequest{
			Table:    "menu_topings",
			ImageUrl: mt.ImageURL,
		})

		if err != nil {
			uc.log.Error(fmt.Errorf("failed to delete image: %w", err).Error())
		}
	}()

	return nil
}
