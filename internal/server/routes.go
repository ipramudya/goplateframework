package server

import (
	"github.com/goplateframework/internal/domain/account/accountdelivery"
	"github.com/goplateframework/internal/domain/account/accountrepo"
	"github.com/goplateframework/internal/domain/account/accountuc"
	"github.com/goplateframework/internal/domain/auth/authdelivery"
	"github.com/goplateframework/internal/domain/auth/authrepo"
	"github.com/goplateframework/internal/domain/auth/authuc"
	"github.com/goplateframework/internal/domain/outlet/outletdelivery"
	"github.com/goplateframework/internal/domain/outlet/outletrepo"
	"github.com/goplateframework/internal/domain/outlet/outletuc"
	"github.com/goplateframework/internal/web"
)

func routes(w *web.Web, conf *Config) {
	authCacheRepo := authrepo.NewCache(conf.RDB)
	accountDBRepo := accountrepo.NewDB(conf.DB)
	outletDBRepo := outletrepo.NewDB(conf.DB)

	accountUC := accountuc.New(conf.ServConf, conf.Log, accountDBRepo)
	authUC := authuc.New(conf.ServConf, conf.Log, authCacheRepo, accountDBRepo)
	outletUC := outletuc.New(conf.ServConf, conf.Log, outletDBRepo)

	accountdelivery.Route(w, &accountdelivery.Options{
		Log:       conf.Log,
		AccountUC: accountUC,
	})

	authdelivery.Route(w, &authdelivery.Options{
		Log:    conf.Log,
		AuthUC: authUC,
	})

	outletdelivery.Route(w, &outletdelivery.Options{
		Log:      conf.Log,
		OutletUC: outletUC,
	})
}
