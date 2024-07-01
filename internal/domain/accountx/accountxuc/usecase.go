package accountxuc

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/domain/accountx"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type iDBRepository interface {
	GetOne(ctx context.Context, id uuid.UUID) (*accountx.AccountDTO, error)
	GetOneByEmail(ctx context.Context, email string) (*accountx.AccountDTO, error)
	Create(ctx context.Context, a *accountx.AccountDTO) error
	ChangePassword(ctx context.Context, email, password string) error
}

type iCacheRepository interface {
	SetMe(ctx context.Context, accountPayload *accountx.AccountDTO) error
	GetMe(ctx context.Context, id uuid.UUID) (*accountx.AccountDTO, error)
}

type Usecase struct {
	conf      *config.Config
	log       *logger.Log
	dbRepo    iDBRepository
	cacheRepo iCacheRepository
}

func New(conf *config.Config, log *logger.Log, dbRepo iDBRepository, cacheRepo iCacheRepository) *Usecase {
	return &Usecase{
		conf:      conf,
		log:       log,
		dbRepo:    dbRepo,
		cacheRepo: cacheRepo,
	}
}

func (uc *Usecase) Register(ctx context.Context, na *accountx.NewAccouuntDTO) (*accountx.AccountDTO, error) {
	existing, err := uc.dbRepo.GetOneByEmail(ctx, na.Email)
	if existing != nil && err == nil {
		e := errs.Newf(errs.AlreadyExists, "email %s already exists", na.Email)
		uc.log.Error(e.Debug())
		return nil, e
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(na.Password), bcrypt.DefaultCost)
	if err != nil {
		e := errs.Newf(errs.Internal, "failed to hash password: %v", err)
		uc.log.Error(e.Debug())
		return nil, e
	}

	now := time.Now()
	a := &accountx.AccountDTO{
		ID:        uuid.New(),
		Firstname: na.Firstname,
		Lastname:  na.Lastname,
		Phone:     na.Phone,
		Email:     na.Email,
		Password:  string(hashedPass),
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.dbRepo.Create(ctx, a); err != nil {
		e := errs.Newf(errs.Internal, "failed to create account: %v", err)
		uc.log.Error(e.Debug())
		return nil, e
	}

	return a, nil
}

func (uc *Usecase) ChangePassword(ctx context.Context, cp *accountx.ChangePasswordDTO, email string) error {
	a, err := uc.dbRepo.GetOneByEmail(ctx, email)
	if err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.Debug())
		return e
	}

	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(cp.OldPassword)); err != nil {
		e := errs.New(errs.InvalidCredentials, errors.New("invalid email or password"))
		uc.log.Error(e.Debug())
		return e
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(cp.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		e := errs.Newf(errs.Internal, "failed to hash password: %v", err)
		uc.log.Error(e.Debug())
		return e
	}

	if err := uc.dbRepo.ChangePassword(ctx, email, string(hashedPass)); err != nil {
		e := errs.Newf(errs.Internal, "failed to change password: %v", err)
		uc.log.Error(e.Debug())
		return e
	}

	return nil
}

func (uc *Usecase) Me(ctx context.Context, accountID uuid.UUID) (*accountx.AccountDTO, error) {
	meCache, err := uc.cacheRepo.GetMe(ctx, accountID)
	if err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.Debug())
		return nil, e
	}

	if meCache != nil {
		return meCache, nil
	}

	a, err := uc.dbRepo.GetOne(ctx, accountID)

	if err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.Debug())
		return nil, e
	}

	if err := uc.cacheRepo.SetMe(ctx, a); err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.Debug())
		return nil, e

	}

	return a, nil
}
