package middlewares

import (
	"errors"

	"github.com/goplateframework/internal/sdk/errs"
	"github.com/goplateframework/internal/sdk/jsonwebtoken"
	"github.com/goplateframework/internal/web/webcontext"
	"github.com/labstack/echo/v4"
)

func (mid *Middleware) Authenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			e := errs.Newf(errs.Unauthenticated, "unauthenticated")
			mid.log.Debug(e.Debug())
			return c.JSON(e.HTTPStatus(), e)
		}

		token, err := jsonwebtoken.ExtractBearerToken(authHeader)
		if err != nil {
			e := errs.Newf(errs.Unauthenticated, "unauthenticated: %v", err)
			mid.log.Debug(e.Debug())
			return c.JSON(e.HTTPStatus(), e)
		}

		// whenever token exist in blacklist, it will return error
		val, _ := mid.cache.Get(c.Request().Context(), token).Result()
		if val != "" {
			e := errs.New(errs.Unauthenticated, errors.New("unauthenticated: expired token"))
			mid.log.Debug(e.Debug())
			return c.JSON(e.HTTPStatus(), e)
		}

		claims, err := jsonwebtoken.Validate(mid.conf, token)
		if err != nil {
			if errors.Is(err, jsonwebtoken.ErrInvalidToken) {
				e := errs.Newf(errs.Unauthenticated, "unauthenticated: %v", err)
				mid.log.Debug(e.Debug())
				return c.JSON(e.HTTPStatus(), e)
			}

			e := errs.New(errs.Internal, err)
			mid.log.Debug(e.Debug())
			return c.JSON(e.HTTPStatus(), e)
		}

		ctx := webcontext.SetClaims(c.Request().Context(), claims)
		ctx = webcontext.SetToken(ctx, token)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
