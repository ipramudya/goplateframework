package authweb

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/goplateframework/internal/domain/auth"
	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/internal/sdk/tokenutil"
	"github.com/goplateframework/internal/web/webcontext"
	"github.com/goplateframework/pkg/logger"
	"github.com/labstack/echo/v4"
)

type iUsecase interface {
	Login(ctx context.Context, email, password string) (*auth.AuthDTO, error)
	Logout(ctx context.Context, accessToken, refreshToken string, atc *tokenutil.AccessTokenClaims, rtc *tokenutil.RefreshTokenClaims) error
	Refresh(ctx context.Context, refreshToken string, accountID uuid.UUID) (*auth.AuthDTO, error)
}

type controller struct {
	authUC iUsecase
	log    *logger.Log
}

func newController(authUC iUsecase, log *logger.Log) *controller {
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

	a, err := con.authUC.Login(c.Request().Context(), dto.Email, dto.Password)
	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusOK, a)
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

	a, err := con.authUC.Refresh(c.Request().Context(), rt, claims.AccountID)

	if err != nil {
		return c.JSON(err.(*errs.Error).HTTPStatus(), err)
	}

	return c.JSON(http.StatusOK, a)
}
