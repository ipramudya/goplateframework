package authuc

import (
	"context"
	"errors"
	"sync"

	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/domain/account"
	"github.com/goplateframework/internal/domain/auth"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/internal/sdk/tokenutil"
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

	if err := uc.cache.RemoveRefreshTokenFromBlacklist(ctx, account.ID.String()); err != nil {
		e := errs.Newf(errs.Internal, "failed to remove refresh token from blacklist: %v", err)
		uc.log.Error(e.Debug())
		return &auth.AccountWithTokenDTO{}, e
	}

	atCh, rtCh := make(chan *string, 1), make(chan *string, 1) //access token channel & refresh token channel
	uc.wg.Add(2)

	go func(ch chan *string) {
		defer func() {
			uc.wg.Done()
			close(ch)
		}()
		at, err := tokenutil.GenerateAccess(uc.conf, tokenutil.AccessTokenPayload{
			AccountID: account.ID.String(),
			Email:     account.Email,
			Role:      account.Role,
		})
		if err != nil {
			e := errs.Newf(errs.Internal, "failed to generate access_token: %v", err)
			uc.log.Error(e.Debug())
			ch <- nil // store nil pointer to channel
			return
		}
		ch <- &at
	}(atCh)

	go func(ch chan *string) {
		defer func() {
			uc.wg.Done()
			close(ch)
		}()
		rt, err := tokenutil.GenerateRefresh(uc.conf, tokenutil.RefreshTokenPayload{
			AccountID: account.ID.String(),
		})
		if err != nil {
			e := errs.Newf(errs.Internal, "failed to generate refresh_token: %v", err)
			uc.log.Error(e.Debug())
			ch <- nil // store nil pointer to channel
			return
		}
		ch <- &rt
	}(rtCh)

	uc.wg.Wait()
	at, rt := <-atCh, <-rtCh

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

func (uc *Usecase) Logout(ctx context.Context, accessToken, refreshToken string, atc *tokenutil.AccessTokenClaims, rtc *tokenutil.RefreshTokenClaims) error {
	atRemaining := tokenutil.RemainingTime(&atc.RegisteredClaims)
	rtRemaining := tokenutil.RemainingTime(&rtc.RegisteredClaims)

	chanErrs := make(chan error, 2)
	uc.wg.Add(2)

	go func(e chan error) {
		defer uc.wg.Done()
		if err := uc.cache.AddAccessTokenToBlacklist(ctx, accessToken, atRemaining); err != nil {
			e := errs.Newf(errs.Internal, "failed to add access token to blacklist: %v", err)
			uc.log.Error(e.Debug())
			chanErrs <- e
		}
		chanErrs <- nil
	}(chanErrs)

	go func(e chan error) {
		defer uc.wg.Done()
		if err := uc.cache.AddRefreshTokenToBlacklist(ctx, rtc.AccountID, refreshToken, rtRemaining); err != nil {
			e := errs.Newf(errs.Internal, "failed to add refresh token to blacklist: %v", err)
			uc.log.Error(e.Debug())
			chanErrs <- e
		}
		chanErrs <- nil
	}(chanErrs)

	uc.wg.Wait()
	e := <-chanErrs
	close(chanErrs)

	if e != nil {
		return e
	}

	return nil
}

func (uc *Usecase) Refresh(ctx context.Context, refreshToken, accountID string) (*auth.AccountWithTokenDTO, error) {
	account, err := uc.accountRepo.GetOneByID(ctx, accountID)
	if err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.Debug())
		return &auth.AccountWithTokenDTO{}, e
	}

	if account == nil {
		e := errs.New(errs.NotFound, errors.New("account not found"))
		uc.log.Error(e.Debug())
		return &auth.AccountWithTokenDTO{}, e
	}

	at, err := tokenutil.GenerateAccess(uc.conf, tokenutil.AccessTokenPayload{
		AccountID: account.ID.String(),
		Email:     account.Email,
		Role:      account.Role,
	})

	if err != nil {
		e := errs.Newf(errs.Internal, "failed to generate access_token: %v", err)
		uc.log.Error(e.Debug())
		return &auth.AccountWithTokenDTO{}, e
	}

	return &auth.AccountWithTokenDTO{
		Account:      account.IntoAccountDTO(),
		AccessToken:  at,
		RefreshToken: refreshToken,
	}, nil
}
