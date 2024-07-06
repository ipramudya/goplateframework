package authweb

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/goplateframework/internal/domain/auth"
	"github.com/goplateframework/internal/sdk/errshttp"
	"github.com/goplateframework/internal/sdk/tokenutil"
	"github.com/goplateframework/internal/sdk/validate"
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
		return errshttp.New(errshttp.InvalidArgument, "Given JSON is invalid")
	}

	if err := dto.Validate(); err != nil {
		e := errshttp.New(errshttp.InvalidArgument, "Given JSON is out of validation rules")

		validationErrs := validate.SplitErrors(err)
		for _, s := range validationErrs {
			e.AddDetail(s)
		}

		return e
	}

	a, err := con.authUC.Login(c.Request().Context(), dto.Email, dto.Password)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, a)
}

func (con *controller) logout(c echo.Context) error {
	at := webcontext.GetAccessToken(c.Request().Context())
	rt := webcontext.GetRefreshToken(c.Request().Context())

	if at == "" || rt == "" {
		e := errshttp.New(errshttp.Unauthenticated, "Could not give access to this resource")
		if at == "" {
			e.AddDetail("token: access_token is missing")
		}

		if rt == "" {
			e.AddDetail("token: refresh_token is missing")
		}

		return e
	}

	atc := webcontext.GetAccessTokenClaims(c.Request().Context())
	rtc := webcontext.GetRefreshTokenClaims(c.Request().Context())

	if atc == nil || rtc == nil {
		e := errshttp.New(errshttp.Unauthenticated, "Could not give access to this resource")
		e.AddDetail("data: claims are not found")
		return e
	}

	err := con.authUC.Logout(c.Request().Context(), at, rt, atc, rtc)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (con *controller) refreshToken(c echo.Context) error {
	rt := webcontext.GetRefreshToken(c.Request().Context())

	if rt == "" {
		e := errshttp.New(errshttp.Unauthenticated, "Could not give access to this resource")
		e.AddDetail("token: refresh_token is missing")
		return e
	}

	claims := webcontext.GetRefreshTokenClaims(c.Request().Context())

	if claims == nil {
		e := errshttp.New(errshttp.Unauthenticated, "Could not give access to this resource")
		e.AddDetail("data: claims are not found")
		return e
	}

	a, err := con.authUC.Refresh(c.Request().Context(), rt, claims.AccountID)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, a)
}
