package server

import (
	"github.com/goplateframework/internal/domain/account/accountdelivery"
	"github.com/goplateframework/internal/domain/account/accountrepo"
	"github.com/goplateframework/internal/domain/account/accountuc"
	"github.com/goplateframework/internal/web"
)

func routes(w *web.Web, conf *Config) {
	// instantiate all domain use cases & repositories for cross requirements
	accountRepo := accountrepo.New(conf.DB)
	accountUC := accountuc.New(conf.ServConf, conf.Log, accountRepo)

	// account domain delivery/endpoint
	accountdelivery.Route(w, &accountdelivery.Options{
		Log:       conf.Log,
		AccountUC: accountUC,
	})
}
