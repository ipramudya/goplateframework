package httpserver

import (
	"net/http"

	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/web"
	"github.com/goplateframework/internal/worker/pb"
	"github.com/goplateframework/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Options struct {
	DB       *sqlx.DB
	Cache    *redis.Client
	Log      *logger.Log
	ServConf *config.Config
	Worker   pb.WorkerClient
}

func Handler(conf *Options) http.Handler {
	// create web application which contains a echo instance, as well as http server
	w := web.New(conf.Log)

	// middleware setup
	w.InitCustomMware(conf.ServConf, conf.Cache)
	w.EnableCORSMware(conf.ServConf.Server.AllowedOrigins)
	w.EnableRecovererMware()
	w.EnableGlobalMware()

	// remap all routes to the web application
	router(w, conf)

	return w
}
