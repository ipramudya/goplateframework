package authdelivery

import (
	"errors"
	"net/http"

	"github.com/goplateframework/internal/domain/auth"
	"github.com/goplateframework/internal/domain/auth/authuc"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/internal/web/webcontext"
	"github.com/goplateframework/pkg/logger"
	"github.com/labstack/echo/v4"
)

type controller struct {
	authUC *authuc.Usecase
	log    *logger.Log
}

func newController(authUC *authuc.Usecase, log *logger.Log) *controller {
	return &controller{authUC, log}
}

func (con *controller) login(c echo.Context) error {
	dto := new(auth.LoginDTO)

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

	account, err := con.authUC.Login(c.Request().Context(), dto.Email, dto.Password)
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusOK, account)
}

func (con *controller) logout(c echo.Context) error {
	at := webcontext.GetAccessToken(c.Request().Context())
	rt := webcontext.GetRefreshToken(c.Request().Context())

	if at == "" || rt == "" {
		e := errs.New(errs.Unauthenticated, errors.New("unauthenticated"))
		return e
	}

	atc := webcontext.GetAccessTokenClaims(c.Request().Context())
	rtc := webcontext.GetRefreshTokenClaims(c.Request().Context())

	if atc == nil || rtc == nil {
		e := errs.New(errs.Unauthenticated, errors.New("unauthenticated"))
		con.log.Error(e.Debug())
		return e
	}

	err := con.authUC.Logout(c.Request().Context(), at, rt, atc, rtc)
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (con *controller) refreshToken(c echo.Context) error {
	rt := webcontext.GetRefreshToken(c.Request().Context())

	if rt == "" {
		e := errs.New(errs.Unauthenticated, errors.New("unauthenticated"))
		return c.JSON(e.HTTPStatus(), e)
	}

	claims := webcontext.GetRefreshTokenClaims(c.Request().Context())

	if claims == nil {
		e := errs.New(errs.Unauthenticated, errors.New("unauthenticated"))
		return c.JSON(e.HTTPStatus(), e)
	}

	account, err := con.authUC.Refresh(c.Request().Context(), rt, claims.AccountID)

	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusOK, account)
}
