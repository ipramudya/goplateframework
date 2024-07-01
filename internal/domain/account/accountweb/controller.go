package accountweb

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/goplateframework/internal/domain/account"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/internal/web/webcontext"
	"github.com/goplateframework/pkg/logger"
	"github.com/labstack/echo/v4"
)

type iUsecase interface {
	Register(ctx context.Context, na *account.NewAccouuntDTO) (*account.AccountDTO, error)
	ChangePassword(ctx context.Context, cp *account.ChangePasswordDTO, email string) error
	Me(ctx context.Context, accountID uuid.UUID) (*account.AccountDTO, error)
}

type controller struct {
	accountUC iUsecase
	log       *logger.Log
}

func newController(accountUC iUsecase, log *logger.Log) *controller {
	return &controller{accountUC, log}
}

func (con *controller) register(c echo.Context) error {
	dto := new(account.NewAccouuntDTO)

	if err := c.Bind(dto); err != nil {
		e := errs.Newf(errs.InvalidArgument, "invalid request: %v", err)
		con.log.Error(e.Debug())
		return c.JSON(e.HTTPStatus(), e)
	}

	if err := dto.Validate(); err != nil {
		e := errs.Newf(errs.InvalidArgument, "invalid request: (%v)", err)
		con.log.Error(e.Debug())
		return c.JSON(e.HTTPStatus(), e)
	}

	a, err := con.accountUC.Register(c.Request().Context(), dto)
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusCreated, a)
}

func (con *controller) changePassword(c echo.Context) error {
	claims := webcontext.GetAccessTokenClaims(c.Request().Context())
	if claims == nil {
		e := errs.New(errs.Unauthenticated, errors.New("unauthenticated"))
		con.log.Error(e.Debug())
		return e
	}

	dto := new(account.ChangePasswordDTO)

	if err := c.Bind(dto); err != nil {
		e := errs.Newf(errs.InvalidArgument, "invalid request: %v", err)
		con.log.Error(e.Debug())
		return c.JSON(e.HTTPStatus(), e)
	}

	if err := dto.Validate(); err != nil {
		e := errs.Newf(errs.InvalidArgument, "invalid request: (%v)", err)
		con.log.Error(e.Debug())
		return c.JSON(e.HTTPStatus(), e)
	}

	if err := con.accountUC.ChangePassword(c.Request().Context(), dto, claims.Email); err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (con *controller) me(c echo.Context) error {
	claims := webcontext.GetAccessTokenClaims(c.Request().Context())
	if claims == nil {
		e := errs.New(errs.Unauthenticated, errors.New("unauthenticated"))
		return c.JSON(e.HTTPStatus(), e)
	}

	a, err := con.accountUC.Me(c.Request().Context(), claims.AccountID)
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusOK, a)
}
