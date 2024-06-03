package authuc

import (
	"context"

	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/domain/account"
	"github.com/goplateframework/internal/domain/auth"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/internal/sdk/tokenutil"
	"github.com/goplateframework/internal/web/webcontext"
	"github.com/goplateframework/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type Usecase struct {
	conf        *config.Config
	log         *logger.Log
	cache       auth.CacheRepository
	accountRepo account.DBRepository
}

func New(conf *config.Config, log *logger.Log, cache auth.CacheRepository, repo account.DBRepository) *Usecase {
	return &Usecase{
		conf:        conf,
		log:         log,
		cache:       cache,
		accountRepo: repo,
	}
}

func (uc *Usecase) Login(ctx context.Context, email, password string) (*account.AccountWithTokenDTO, error) {
	existingAccount, err := uc.accountRepo.GetOneByEmail(ctx, email)

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
