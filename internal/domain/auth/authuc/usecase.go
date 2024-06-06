package authuc

import (
	"context"
	"sync"

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
	accountRepo account.DBRepository
	cache       auth.CacheRepository
	conf        *config.Config
	log         *logger.Log
	wg          *sync.WaitGroup
}

func New(conf *config.Config, log *logger.Log, cache auth.CacheRepository, repo account.DBRepository) *Usecase {
	return &Usecase{
		accountRepo: repo,
		cache:       cache,
		conf:        conf,
		log:         log,
		wg:          new(sync.WaitGroup),
	}
}

func (uc *Usecase) Login(ctx context.Context, email, password string) (*auth.AccountWithTokenDTO, error) {
	account, err := uc.accountRepo.GetOneByEmail(ctx, email)

	if err != nil {
		e := errs.Newf(errs.InvalidCredentials, "invalid email or password")
		uc.log.Error(e.Debug())
		return &auth.AccountWithTokenDTO{}, e
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)); err != nil {
		e := errs.Newf(errs.InvalidCredentials, "invalid email or password")
		uc.log.Error(e.Debug())
		return &auth.AccountWithTokenDTO{}, e
	}

	token := make(chan *string, 2)
	uc.wg.Add(2)

	go func(t chan *string) {
		defer uc.wg.Done()
		at, err := tokenutil.GenerateAccess(uc.conf, tokenutil.AccessTokenPayload{
			AccountID: account.ID.String(),
			Email:     account.Email,
			Role:      account.Role,
		})
		if err != nil {
			e := errs.Newf(errs.Internal, "failed to generate access_token: %v", err)
			uc.log.Error(e.Debug())
			t <- nil // store nil pointer to channel
			return
		}
		t <- &at
	}(token)

	go func(t chan *string) {
		defer uc.wg.Done()
		rt, err := tokenutil.GenerateRefresh(uc.conf, tokenutil.RefreshTokenPayload{
			AccountID: account.ID.String(),
		})
		if err != nil {
			e := errs.Newf(errs.Internal, "failed to generate refresh_token: %v", err)
			uc.log.Error(e.Debug())
			t <- nil // store nil pointer to channel
			return
		}
		t <- &rt
	}(token)

	uc.wg.Wait()
	at := <-token // access token
	rt := <-token // refresh token
	close(token)

	if at == nil || rt == nil {
		e := errs.Newf(errs.Internal, "failed to generate token")
		return &auth.AccountWithTokenDTO{}, e
	}

	return &auth.AccountWithTokenDTO{
		Account:      account.IntoAccountDTO(),
		AccessToken:  *at,
		RefreshToken: *rt,
	}, nil
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

	if err := uc.cache.AddAccessTokenToBlacklist(ctx, claims.AccountID, token, remaining); err != nil {
		e := errs.Newf(errs.Internal, "failed to add token to blacklist: %v", err)
		uc.log.Error(e.Debug())
		return e
	}

	return nil
}
