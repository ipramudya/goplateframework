package middlewares

import (
	"net/http"

	"github.com/goplateframework/config"
	"github.com/goplateframework/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Middleware struct {
	conf  *config.Config
	log   *logger.Log
	cache *redis.Client
}

type MiddlewareFunc func(h http.Handler) http.Handler

func New(conf *config.Config, log *logger.Log, cache *redis.Client) *Middleware {
	return &Middleware{
		conf:  conf,
		log:   log,
		cache: cache,
	}
}
