package server

import (
	"github.com/goplateframework/internal/domain/account/accountdelivery"
	"github.com/goplateframework/internal/domain/account/accountrepo"
	"github.com/goplateframework/internal/domain/account/accountuc"
	"github.com/goplateframework/internal/domain/address/addressdelivery"
	"github.com/goplateframework/internal/domain/address/addressrepo"
	"github.com/goplateframework/internal/domain/address/addressuc"
	"github.com/goplateframework/internal/domain/auth/authdelivery"
	"github.com/goplateframework/internal/domain/auth/authrepo"
	"github.com/goplateframework/internal/domain/auth/authuc"
	"github.com/goplateframework/internal/domain/menu/menudelivery"
	"github.com/goplateframework/internal/domain/menu/menurepo"
	"github.com/goplateframework/internal/domain/menu/menuuc"
	"github.com/goplateframework/internal/domain/outlet/outletdelivery"
	"github.com/goplateframework/internal/domain/outlet/outletrepo"
	"github.com/goplateframework/internal/domain/outlet/outletuc"
	"github.com/goplateframework/internal/web"
)

func routes(w *web.Web, conf *Config) {
	authCacheRepo := authrepo.NewCache(conf.RDB)
	accountDBRepo := accountrepo.NewDB(conf.DB)
	accountCacheRepo := accountrepo.NewCache(conf.RDB)
	outletDBRepo := outletrepo.NewDB(conf.DB)
	addressDBRepo := addressrepo.NewDB(conf.DB)
	menuDBRepo := menurepo.NewDB(conf.DB)

	accountUC := accountuc.New(conf.ServConf, conf.Log, accountDBRepo, accountCacheRepo)
	authUC := authuc.New(conf.ServConf, conf.Log, authCacheRepo, accountDBRepo)
	outletUC := outletuc.New(conf.ServConf, conf.Log, outletDBRepo, addressDBRepo)
	addressUC := addressuc.New(conf.ServConf, conf.Log, addressDBRepo)
	menuUC := menuuc.New(conf.ServConf, conf.Log, menuDBRepo)

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

	addressdelivery.Route(w, &addressdelivery.Options{
		Log:       conf.Log,
		AddressUC: addressUC,
	})

	menudelivery.Route(w, &menudelivery.Options{
		Log:    conf.Log,
		MenuUC: menuUC,
	})
}
