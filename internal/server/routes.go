package server

import (
	"github.com/goplateframework/internal/domain/account/accountdelivery"
	"github.com/goplateframework/internal/domain/account/accountrepo"
	"github.com/goplateframework/internal/domain/account/accountuc"
	"github.com/goplateframework/internal/domain/auth/authdelivery"
	"github.com/goplateframework/internal/domain/auth/authrepo"
	"github.com/goplateframework/internal/domain/auth/authuc"
	"github.com/goplateframework/internal/web"
)

func routes(w *web.Web, conf *Config) {
	authCacheRepo := authrepo.NewCache(conf.RDB)
	accountDBRepo := accountrepo.NewDB(conf.DB)

	accountUC := accountuc.New(conf.ServConf, conf.Log, accountDBRepo)
	authUC := authuc.New(conf.ServConf, conf.Log, authCacheRepo, accountDBRepo)

	// account domain delivery/endpoint
	accountdelivery.Route(w, &accountdelivery.Options{
		Log:       conf.Log,
		AccountUC: accountUC,
	})

	// auth domain delivery/endpoint
	authdelivery.Route(w, &authdelivery.Options{
		Log:    conf.Log,
		AuthUC: authUC,
	})
}
