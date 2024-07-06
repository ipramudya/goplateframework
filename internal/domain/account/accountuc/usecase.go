package accountuc

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/domain/account"
	"github.com/goplateframework/internal/sdk/errshttp"
	"github.com/goplateframework/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type iDBRepository interface {
	GetOne(ctx context.Context, id uuid.UUID) (*account.AccountDTO, error)
	GetOneByEmail(ctx context.Context, email string) (*account.AccountDTO, error)
	Create(ctx context.Context, a *account.AccountDTO) error
	ChangePassword(ctx context.Context, email, password string) error
}

type iCacheRepository interface {
	SetMe(ctx context.Context, accountPayload *account.AccountDTO) error
	GetMe(ctx context.Context, id uuid.UUID) (*account.AccountDTO, error)
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

func (uc *Usecase) Register(ctx context.Context, na *account.NewAccouuntDTO) (*account.AccountDTO, error) {
	existing, err := uc.dbRepo.GetOneByEmail(ctx, na.Email)
	if existing != nil && err == nil {
		e := errshttp.New(errshttp.AlreadyExists, "Email already taken")
		e.AddDetail(fmt.Sprintf("data: email %s already taken", na.Email))
		return nil, e
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(na.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errshttp.New(errshttp.Internal, "Failed to hash password")
	}

	now := time.Now()
	a := &account.AccountDTO{
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
		return nil, errshttp.New(errshttp.Internal, "Something went wrong")
	}

	return a, nil
}

func (uc *Usecase) ChangePassword(ctx context.Context, cp *account.ChangePasswordDTO, email string) error {
	a, err := uc.dbRepo.GetOneByEmail(ctx, email)
	if err != nil {
		errshttp.New(errshttp.Internal, "Something went wrong")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(cp.OldPassword)); err != nil {
		return errshttp.New(errshttp.InvalidCredentials, "Credentials are invalid, either email and/or password")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(cp.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errshttp.New(errshttp.Internal, "Failed to perform password hashing")
	}

	if err := uc.dbRepo.ChangePassword(ctx, email, string(hashedPass)); err != nil {
		errshttp.New(errshttp.Internal, "Something went wrong")
	}

	return nil
}

func (uc *Usecase) Me(ctx context.Context, accountID uuid.UUID) (*account.AccountDTO, error) {
	meCache, err := uc.cacheRepo.GetMe(ctx, accountID)
	if err != nil {
		return nil, errshttp.New(errshttp.Internal, "Could not retrieve data from cache")
	}

	if meCache != nil {
		return meCache, nil
	}

	a, err := uc.dbRepo.GetOne(ctx, accountID)

	if err != nil {
		return nil, errshttp.New(errshttp.Internal, "Something went wrong")
	}

	if err := uc.cacheRepo.SetMe(ctx, a); err != nil {
		return nil, errshttp.New(errshttp.Internal, "Something went wrong")

	}

	return a, nil
}
