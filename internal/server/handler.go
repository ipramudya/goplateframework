package server

import (
	"net/http"

	"github.com/goplateframework/config"
	"github.com/goplateframework/internal/web"
	"github.com/goplateframework/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	DB       *sqlx.DB
	Log      *logger.Log
	ServConf *config.Config
}

func Handler(conf *Config) http.Handler {
	// create web application which contains a echo instance, as well as http server
	w := web.New(conf.Log)

	// middleware setup
	w.InitCustomMware(conf.ServConf)
	w.EnableCORSMware(conf.ServConf.Server.AllowedOrigins)
	w.EnableRecovererMware()
	w.EnableGlobalMware()

	// remap all routes to the web application
	routes(w, conf)

	return w
}
