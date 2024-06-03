package accountuc

import (
	"context"

	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/domain/account"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/internal/sdk/jsonwebtoken"
	"github.com/goplateframework/internal/web/webcontext"
	"github.com/goplateframework/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type Usecase struct {
	conf *config.Config
	log  *logger.Log
	repo account.Repository
}

func New(conf *config.Config, log *logger.Log, repo account.Repository) *Usecase {
	return &Usecase{
		conf: conf,
		log:  log,
		repo: repo,
	}
}

func (uc *Usecase) Register(ctx context.Context, na *account.NewAccouuntDTO) (*account.AccountWithTokenDTO, error) {
	existingAccount, err := uc.repo.GetOneByEmail(ctx, na.Email)
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

	accountCreated, err := uc.repo.Register(ctx, na)
	if err != nil {
		e := errs.Newf(errs.Internal, "failed to create account: %v", err)
		uc.log.Error(e.Debug())
		return &account.AccountWithTokenDTO{}, e
	}

	token, err := jsonwebtoken.Generate(uc.conf, jsonwebtoken.Payload{
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
	existingAccount, err := uc.repo.GetOneByEmail(ctx, email)

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

	token, err := jsonwebtoken.Generate(uc.conf, jsonwebtoken.Payload{
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
	p, err := webcontext.GetAccountPayload(ctx)

	if err != nil {
		e := errs.New(errs.Unauthenticated, err)
		uc.log.Error(e.Debug())
		return e
	}

	account, err := uc.repo.GetOneByEmail(ctx, p.Email)
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

	if err := uc.repo.ChangePassword(ctx, account.Email, string(passHash)); err != nil {
		e := errs.Newf(errs.Internal, "failed to change password: %v", err)
		uc.log.Error(e.Debug())
		return e
	}

	return nil
}
