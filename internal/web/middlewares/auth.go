package middlewares

import (
	"errors"

	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/internal/sdk/tokenutil"
	"github.com/goplateframework/internal/web/webcontext"
	"github.com/labstack/echo/v4"
)

func (mid *Middleware) Authenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")

		token, err := tokenutil.ExtractBearerToken(authHeader)
		if err != nil {
			e := errs.Newf(errs.Unauthenticated, "unauthenticated: %v", err)
			mid.log.Error(e.Debug())
			return c.JSON(e.HTTPStatus(), e)
		}

		// whenever token exist in blacklist, it will return error
		if val, _ := mid.cache.Get(c.Request().Context(), token).Result(); val != "" {
			e := errs.New(errs.Unauthenticated, errors.New("unauthenticated: already logged out"))
			mid.log.Error(e.Debug())
			return c.JSON(e.HTTPStatus(), e)
		}

		claims, err := tokenutil.ValidateAccess(mid.conf, token)
		if err != nil {
			if errors.Is(err, tokenutil.ErrInvalidToken) {
				e := errs.Newf(errs.Unauthenticated, "unauthenticated: %v", err)
				mid.log.Error(e.Debug())
				return c.JSON(e.HTTPStatus(), e)
			}

			e := errs.New(errs.Internal, err)
			mid.log.Error(e.Debug())
			return c.JSON(e.HTTPStatus(), e)
		}

		ctx := webcontext.SetAccessTokenClaims(c.Request().Context(), claims)
		ctx = webcontext.SetAccessToken(ctx, token)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

func (mid *Middleware) RefreshAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		refreshToken := c.Request().Header.Get("RF-Token")
		if refreshToken == "" {
			e := errs.New(errs.Unauthenticated, errors.New("unauthenticated"))
			mid.log.Error(e.Debug())
			return c.JSON(e.HTTPStatus(), e)
		}

		claims, err := tokenutil.ValidateRefresh(mid.conf, refreshToken)
		if err != nil {
			e := errs.New(errs.Unauthenticated, errors.New("unauthenticated"))
			mid.log.Error(e.Debug())
			return c.JSON(e.HTTPStatus(), e)
		}

		// whenever token exist in blacklist, it will return error
		if val, _ := mid.cache.Get(c.Request().Context(), claims.AccountID).Result(); val != "" {
			e := errs.New(errs.Unauthenticated, errors.New("unauthenticated: already logged out"))
			mid.log.Error(e.Debug())
			return c.JSON(e.HTTPStatus(), e)
		}

		ctx := webcontext.SetRefreshTokenClaims(c.Request().Context(), claims)
		ctx = webcontext.SetRefreshToken(ctx, refreshToken)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
