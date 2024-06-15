package accountuc

import (
	"context"
	"errors"

	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/domain/account"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/internal/web/webcontext"
	"github.com/goplateframework/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type Usecase struct {
	conf             *config.Config
	log              *logger.Log
	accountDBRepo    account.DBRepository
	accountCacheRepo account.CacheRepository
}

func New(conf *config.Config, log *logger.Log, accountDBRepo account.DBRepository, accountCacheRepo account.CacheRepository) *Usecase {
	return &Usecase{
		conf:             conf,
		log:              log,
		accountDBRepo:    accountDBRepo,
		accountCacheRepo: accountCacheRepo,
	}
}

func (uc *Usecase) Register(ctx context.Context, na *account.NewAccouuntDTO) (*account.AccountDTO, error) {
	existingAccount, err := uc.accountDBRepo.GetOneByEmail(ctx, na.Email)
	if existingAccount != nil && err == nil {
		e := errs.Newf(errs.AlreadyExists, "email %s already exists", na.Email)
		uc.log.Error(e.Debug())
		return &account.AccountDTO{}, e
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(na.Password), bcrypt.DefaultCost)
	if err != nil {
		e := errs.Newf(errs.Internal, "failed to hash password: %v", err)
		uc.log.Error(e.Debug())
		return &account.AccountDTO{}, e
	}
	na.Password = string(passHash)

	accountCreated, err := uc.accountDBRepo.Register(ctx, na)
	if err != nil {
		e := errs.Newf(errs.Internal, "failed to create account: %v", err)
		uc.log.Error(e.Debug())
		return &account.AccountDTO{}, e
	}

	return accountCreated.IntoAccountDTO(), nil
}

func (uc *Usecase) ChangePassword(ctx context.Context, oldpass, newpass string) error {
	claims := webcontext.GetAccessTokenClaims(ctx)

	if claims == nil {
		e := errs.New(errs.Unauthenticated, errors.New("unauthenticated"))
		uc.log.Error(e.Debug())
		return e
	}

	account, err := uc.accountDBRepo.GetOneByEmail(ctx, claims.Email)
	if err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.Debug())
		return e
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(oldpass)); err != nil {
		e := errs.New(errs.InvalidCredentials, errors.New("invalid email or password"))
		uc.log.Error(e.Debug())
		return e
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(newpass), bcrypt.DefaultCost)
	if err != nil {
		e := errs.Newf(errs.Internal, "failed to hash password: %v", err)
		uc.log.Error(e.Debug())
		return e
	}

	if err := uc.accountDBRepo.ChangePassword(ctx, account.Email, string(passHash)); err != nil {
		e := errs.Newf(errs.Internal, "failed to change password: %v", err)
		uc.log.Error(e.Debug())
		return e
	}

	return nil
}

func (uc *Usecase) Me(ctx context.Context, accountID string) (*account.AccountDTO, error) {
	meCache, err := uc.accountCacheRepo.GetMe(ctx, accountID)

	if err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.Debug())
		return &account.AccountDTO{}, e
	}

	if meCache == nil {
		a, err := uc.accountDBRepo.GetOneByID(ctx, accountID)

		if err != nil {
			e := errs.New(errs.Internal, errors.New("something went wrong"))
			uc.log.Error(e.Debug())
			return &account.AccountDTO{}, e
		}

		if err := uc.accountCacheRepo.SetMe(ctx, a); err != nil {
			e := errs.New(errs.Internal, errors.New("something went wrong"))
			uc.log.Error(e.Debug())
			return &account.AccountDTO{}, e
		}

		return a.IntoAccountDTO(), nil
	} else {
		return meCache.IntoAccountDTO(), nil
	}
}
