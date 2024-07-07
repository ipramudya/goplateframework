package httpserver

import (
	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/web"
	"github.com/goplateframework/internal/worker/pb"
	"github.com/goplateframework/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type Options struct {
	DB       *sqlx.DB
	Cache    *redis.Client
	Log      *logger.Log
	ServConf *config.Config
	Worker   pb.WorkerClient
}

func Init(opts *Options) *echo.Echo {
	// create web application which contains a echo instance, as well as http server
	w := web.New(opts.Log)
	w.Echo.HideBanner = true
	w.Echo.HidePort = true

	// middleware setup
	w.InitCustomMware(opts.ServConf, opts.Cache)
	w.EnableCORSMware(opts.ServConf.Server.AllowedOrigins)
	w.EnableRecovererMware()
	w.EnableGlobalMware()

	// remap all routes to the web application
	router(w, opts)

	return w.Echo
}
