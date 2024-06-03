package middlewares

import (
	"net/http"

	"github.com/goplateframework/config"
	"github.com/goplateframework/pkg/logger"
)

type Middleware struct {
	conf *config.Config
	log  *logger.Log
}

type MiddlewareFunc func(h http.Handler) http.Handler

func New(conf *config.Config, log *logger.Log) *Middleware {
	return &Middleware{conf, log}
}
