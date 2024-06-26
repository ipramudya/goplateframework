package accountweb

import (
	"errors"
	"net/http"

	"github.com/goplateframework/internal/domain/account"
	"github.com/goplateframework/internal/domain/account/accountuc"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/internal/web/webcontext"
	"github.com/goplateframework/pkg/logger"
	"github.com/labstack/echo/v4"
)

type controller struct {
	accountUC *accountuc.Usecase
	log       *logger.Log
}

func newController(accountUC *accountuc.Usecase, log *logger.Log) *controller {
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

	account, err := con.accountUC.Register(c.Request().Context(), dto)
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusCreated, account)
}

func (con *controller) changePassword(c echo.Context) error {
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

	err := con.accountUC.ChangePassword(c.Request().Context(), dto.OldPassword, dto.NewPassword)
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.NoContent(http.StatusAccepted)
}

func (con *controller) me(c echo.Context) error {
	claims := webcontext.GetAccessTokenClaims(c.Request().Context())

	if claims == nil {
		e := errs.New(errs.Unauthenticated, errors.New("unauthenticated"))
		return c.JSON(e.HTTPStatus(), e)
	}

	account, err := con.accountUC.Me(c.Request().Context(), claims.AccountID)
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusOK, account)
}
