package middlewares

import (
	"errors"
	"fmt"

	"github.com/goplateframework/internal/sdk/errshttp"
	"github.com/goplateframework/internal/sdk/tokenutil"
	"github.com/goplateframework/internal/web/webcontext"
	"github.com/labstack/echo/v4"
)

func (mid *Middleware) Authenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")

		token, err := tokenutil.ExtractBearerToken(authHeader)
		if err != nil {
			e := errshttp.New(errshttp.Unauthenticated, "Authorization header missing")
			e.AddDetail(fmt.Sprintf("data: %s", err))
			return e
		}

		// whenever token exist in blacklist, it will return error
		if val, _ := mid.cache.Get(c.Request().Context(), token).Result(); val != "" {
			e := errshttp.New(errshttp.Unauthenticated, "User already logged out")
			e.AddDetail(fmt.Sprintf("data: %s", err))
			return e
		}

		claims, err := tokenutil.ValidateAccess(mid.conf, token)
		if err != nil {
			if errors.Is(err, tokenutil.ErrInvalidToken) {
				e := errshttp.New(errshttp.Unauthenticated, "Invalid access token")
				e.AddDetail("token: access_token on bearer authentication header is invalid")
				return e
			}

			return errshttp.New(errshttp.Internal, "Something went wrong")
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
			e := errshttp.New(errshttp.Unauthenticated, "RT-Token header missing")
			e.AddDetail("token: expected RT-Token header with refresh token as its value")
		}

		claims, err := tokenutil.ValidateRefresh(mid.conf, refreshToken)
		if err != nil {
			e := errshttp.New(errshttp.Unauthenticated, "Invalid refresh token")
			e.AddDetail("token: refresh_token on RT-Token header is invalid")
			return e
		}

		// whenever token exist in blacklist, it will return error
		if val, _ := mid.cache.Get(c.Request().Context(), claims.AccountID.String()).Result(); val != "" {
			e := errshttp.New(errshttp.Unauthenticated, "User already logged out")
			e.AddDetail(fmt.Sprintf("data: %s", err))
			return e
		}

		ctx := webcontext.SetRefreshTokenClaims(c.Request().Context(), claims)
		ctx = webcontext.SetRefreshToken(ctx, refreshToken)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
