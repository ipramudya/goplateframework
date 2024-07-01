package authuc

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/domain/account"
	"github.com/goplateframework/internal/domain/auth"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/internal/sdk/tokenutil"
	"github.com/goplateframework/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type iAuthCacheRepo interface {
	AddAccessTokenToBlacklist(ctx context.Context, token string, exp time.Duration) error
	AddRefreshTokenToBlacklist(ctx context.Context, accountID uuid.UUID, token string, exp time.Duration) error
	RemoveRefreshTokenFromBlacklist(ctx context.Context, accountID uuid.UUID) error
}

type iAccountDBRepo interface {
	GetOneByEmail(ctx context.Context, email string) (*account.AccountDTO, error)
	GetOne(ctx context.Context, id uuid.UUID) (*account.AccountDTO, error)
}

type Usecase struct {
	authCacheRepo iAuthCacheRepo
	conf          *config.Config
	log           *logger.Log
	accountDBRepo iAccountDBRepo
}

func New(conf *config.Config, log *logger.Log, authCacheRepo iAuthCacheRepo, accountDBRepo iAccountDBRepo) *Usecase {
	return &Usecase{
		accountDBRepo: accountDBRepo,
		authCacheRepo: authCacheRepo,
		conf:          conf,
		log:           log,
	}
}

func (uc *Usecase) Login(ctx context.Context, email, password string) (*auth.AuthDTO, error) {
	a, err := uc.accountDBRepo.GetOneByEmail(ctx, email)

	if err != nil {
		e := errs.New(errs.InvalidCredentials, errors.New("invalid email or password"))
		uc.log.Error(e.Debug())
		return nil, e
	}

	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password)); err != nil {
		e := errs.New(errs.InvalidCredentials, errors.New("invalid email or password"))
		uc.log.Error(e.Debug())
		return nil, e
	}

	if err := uc.authCacheRepo.RemoveRefreshTokenFromBlacklist(ctx, a.ID); err != nil {
		e := errs.Newf(errs.Internal, "failed to remove refresh token from blacklist: %v", err)
		uc.log.Error(e.Debug())
		return nil, e
	}

	wg := new(sync.WaitGroup)
	atCh, rtCh := make(chan *string, 1), make(chan *string, 1) //access token channel & refresh token channel
	wg.Add(2)

	go func(ch chan *string) {
		defer wg.Done()

		at, err := tokenutil.GenerateAccess(uc.conf, tokenutil.AccessTokenPayload{
			AccountID: a.ID,
			Email:     a.Email,
			Role:      a.Role,
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
		defer wg.Done()

		rt, err := tokenutil.GenerateRefresh(uc.conf, tokenutil.RefreshTokenPayload{
			AccountID: a.ID,
		})
		if err != nil {
			e := errs.Newf(errs.Internal, "failed to generate refresh_token: %v", err)
			uc.log.Error(e.Debug())
			ch <- nil // store nil pointer to channel
			return
		}
		ch <- &rt
	}(rtCh)

	wg.Wait()
	at, rt := <-atCh, <-rtCh

	if at == nil || rt == nil {
		e := errs.New(errs.Internal, errors.New("failed to generate token"))
		return nil, e
	}

	return &auth.AuthDTO{
		Account:      a,
		AccessToken:  *at,
		RefreshToken: *rt,
	}, nil
}

func (uc *Usecase) Logout(ctx context.Context, accessToken, refreshToken string, atc *tokenutil.AccessTokenClaims, rtc *tokenutil.RefreshTokenClaims) error {
	atTime := tokenutil.RemainingTime(&atc.RegisteredClaims)
	rtTime := tokenutil.RemainingTime(&rtc.RegisteredClaims)

	chanErrs := make(chan error, 2)
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func(e chan error) {
		defer wg.Done()

		err := uc.authCacheRepo.AddAccessTokenToBlacklist(ctx, accessToken, atTime)
		if err != nil {
			e := errs.Newf(errs.Internal, "failed to add access token to blacklist: %v", err)
			uc.log.Error(e.Debug())
			chanErrs <- e
		}
		chanErrs <- nil
	}(chanErrs)

	go func(e chan error) {
		defer wg.Done()

		err := uc.authCacheRepo.AddRefreshTokenToBlacklist(ctx, rtc.AccountID, refreshToken, rtTime)
		if err != nil {
			e := errs.Newf(errs.Internal, "failed to add refresh token to blacklist: %v", err)
			uc.log.Error(e.Debug())
			chanErrs <- e
		}
		chanErrs <- nil
	}(chanErrs)

	wg.Wait()
	e := <-chanErrs

	if e != nil {
		return e
	}

	return nil
}

func (uc *Usecase) Refresh(ctx context.Context, refreshToken string, accountID uuid.UUID) (*auth.AuthDTO, error) {
	a, err := uc.accountDBRepo.GetOne(ctx, accountID)
	if err != nil {
		e := errs.New(errs.Internal, errors.New("something went wrong"))
		uc.log.Error(e.Debug())
		return nil, e
	}

	if a == nil {
		e := errs.New(errs.NotFound, errors.New("account not found"))
		uc.log.Error(e.Debug())
		return nil, e
	}

	at, err := tokenutil.GenerateAccess(uc.conf, tokenutil.AccessTokenPayload{
		AccountID: a.ID,
		Email:     a.Email,
		Role:      a.Role,
	})

	if err != nil {
		e := errs.Newf(errs.Internal, "failed to generate access_token: %v", err)
		uc.log.Error(e.Debug())
		return nil, e
	}

	return &auth.AuthDTO{
		Account:      a,
		AccessToken:  at,
		RefreshToken: refreshToken,
	}, nil
}
