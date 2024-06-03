package accountuc

import (
	"context"
	"errors"

	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/domain/account"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/internal/sdk/tokenutil"
	"github.com/goplateframework/internal/web/webcontext"
	"github.com/goplateframework/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type Usecase struct {
	conf   *config.Config
	log    *logger.Log
	repoDB account.DBRepository
	cache  account.CacheRepository
}

func New(conf *config.Config, log *logger.Log, repo account.DBRepository, cache account.CacheRepository) *Usecase {
	return &Usecase{
		conf:   conf,
		log:    log,
		repoDB: repo,
		cache:  cache,
	}
}

func (uc *Usecase) Register(ctx context.Context, na *account.NewAccouuntDTO) (*account.AccountWithTokenDTO, error) {
	existingAccount, err := uc.repoDB.GetOneByEmail(ctx, na.Email)
	if existingAccount != nil && err == nil {
		e := errs.Newf(errs.AlreadyExists, "email %s already exists", na.Email)
		uc.log.Error(e.Debug())
		return &account.AccountWithTokenDTO{}, e
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(na.Password), bcrypt.DefaultCost)
	if err != nil {
		e := errs.Newf(errs.Internal, "failed to hash password: %v", err)
		uc.log.Error(e.Debug())
		return &account.AccountWithTokenDTO{}, e
	}
	na.Password = string(passHash)

	accountCreated, err := uc.repoDB.Register(ctx, na)
	if err != nil {
		e := errs.Newf(errs.Internal, "failed to create account: %v", err)
		uc.log.Error(e.Debug())
		return &account.AccountWithTokenDTO{}, e
	}

	token, err := tokenutil.Generate(uc.conf, tokenutil.Payload{
		Email:     accountCreated.Email,
		AccountID: accountCreated.ID.String(),
		Role:      accountCreated.Role,
	})

	if err != nil {
		e := errs.Newf(errs.Internal, "failed to generate token: %v", err)
		uc.log.Error(e.Debug())
		return &account.AccountWithTokenDTO{}, e
	}

	return accountCreated.IntoAccountWithTokenDTO(token), nil
}

func (uc *Usecase) Login(ctx context.Context, email, password string) (*account.AccountWithTokenDTO, error) {
	existingAccount, err := uc.repoDB.GetOneByEmail(ctx, email)

	if err != nil {
		e := errs.Newf(errs.InvalidCredentials, "invalid email or password")
		uc.log.Error(e.Debug())
		return &account.AccountWithTokenDTO{}, e
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingAccount.Password), []byte(password)); err != nil {
		e := errs.Newf(errs.InvalidCredentials, "invalid email or password")
		uc.log.Error(e.Debug())
		return &account.AccountWithTokenDTO{}, e
	}

	token, err := tokenutil.Generate(uc.conf, tokenutil.Payload{
		Email:     existingAccount.Email,
		AccountID: existingAccount.ID.String(),
		Role:      existingAccount.Role,
	})

	if err != nil {
		e := errs.Newf(errs.Internal, "failed to generate token: %v", err)
		uc.log.Error(e.Debug())
		return &account.AccountWithTokenDTO{}, e
	}

	return existingAccount.IntoAccountWithTokenDTO(token), nil
}

func (uc *Usecase) ChangePassword(ctx context.Context, oldpass, newpass string) error {
	claims := webcontext.GetClaims(ctx)

	if claims == nil {
		e := errs.New(errs.Unauthenticated, errors.New("unauthenticated"))
		uc.log.Error(e.Debug())
		return e
	}

	account, err := uc.repoDB.GetOneByEmail(ctx, claims.Email)
	if err != nil {
		e := errs.Newf(errs.Internal, "something went wrong!")
		uc.log.Error(e.Debug())
		return e
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(oldpass)); err != nil {
		e := errs.Newf(errs.InvalidCredentials, "invalid email or password")
		uc.log.Error(e.Debug())
		return e
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(newpass), bcrypt.DefaultCost)
	if err != nil {
		e := errs.Newf(errs.Internal, "failed to hash password: %v", err)
		uc.log.Error(e.Debug())
		return e
	}

	if err := uc.repoDB.ChangePassword(ctx, account.Email, string(passHash)); err != nil {
		e := errs.Newf(errs.Internal, "failed to change password: %v", err)
		uc.log.Error(e.Debug())
		return e
	}

	return nil
}

func (uc *Usecase) Logout(ctx context.Context) error {
	token := webcontext.GetToken(ctx)

	if token == "" {
		e := errs.Newf(errs.Unauthenticated, "unauthenticated")
		uc.log.Error(e.Debug())
		return e
	}

	claims := webcontext.GetClaims(ctx)

	if claims == nil {
		e := errs.Newf(errs.Unauthenticated, "unauthenticated")
		uc.log.Error(e.Debug())
		return e
	}

	remaining := tokenutil.RemainingTime(&claims.RegisteredClaims)

	err := uc.cache.AddTokenToBlacklist(ctx, token, remaining)

	if err != nil {
		e := errs.Newf(errs.Internal, "failed to add token to blacklist: %v", err)
		uc.log.Error(e.Debug())
		return e
	}

	return nil
}
