package middlewares

import (
	"time"

	"github.com/labstack/echo/v4"
)

func (mid *Middleware) RequestLoggerMware(h echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		start := time.Now()
		next := h(ctx)

		req := ctx.Request()
		res := ctx.Response()

		since := time.Since(start).String()
		reqId := res.Header().Get(echo.HeaderXRequestID)

		mid.log.Infof("Req ID: %s, Method: %s, URI: %s, Status: %v, Size: %v, Time: %s",
			reqId,
			req.Method,
			req.RequestURI,
			res.Status,
			res.Size,
			since,
		)

		return next
	}
}
