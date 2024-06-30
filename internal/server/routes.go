package server

import (
	"github.com/goplateframework/internal/domain/account/accountrepo"
	"github.com/goplateframework/internal/domain/account/accountuc"
	"github.com/goplateframework/internal/domain/account/accountweb"
	"github.com/goplateframework/internal/domain/address/addressrepo"
	"github.com/goplateframework/internal/domain/address/addressuc"
	"github.com/goplateframework/internal/domain/address/addressweb"
	"github.com/goplateframework/internal/domain/addressx/addressxrepo"
	"github.com/goplateframework/internal/domain/addressx/addressxuc"
	"github.com/goplateframework/internal/domain/addressx/addressxweb"
	"github.com/goplateframework/internal/domain/auth/authrepo"
	"github.com/goplateframework/internal/domain/auth/authuc"
	"github.com/goplateframework/internal/domain/auth/authweb"
	"github.com/goplateframework/internal/domain/menu/menurepo"
	"github.com/goplateframework/internal/domain/menu/menuuc"
	"github.com/goplateframework/internal/domain/menu/menuweb"
	"github.com/goplateframework/internal/domain/menux/menuxrepo"
	"github.com/goplateframework/internal/domain/menux/menuxuc"
	"github.com/goplateframework/internal/domain/menux/menuxweb"
	"github.com/goplateframework/internal/domain/outlet/outletrepo"
	"github.com/goplateframework/internal/domain/outlet/outletuc"
	"github.com/goplateframework/internal/domain/outlet/outletweb"
	"github.com/goplateframework/internal/domain/outletx/outletxrepo"
	"github.com/goplateframework/internal/domain/outletx/outletxuc"
	"github.com/goplateframework/internal/domain/outletx/outletxweb"
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

	accountweb.Route(w, &accountweb.Options{
		Log:       conf.Log,
		AccountUC: accountUC,
	})

	authweb.Route(w, &authweb.Options{
		Log:    conf.Log,
		AuthUC: authUC,
	})

	outletweb.Route(w, &outletweb.Options{
		Log:      conf.Log,
		OutletUC: outletUC,
	})

	addressweb.Route(w, &addressweb.Options{
		Log:       conf.Log,
		AddressUC: addressUC,
	})

	menuweb.Route(w, &menuweb.Options{
		Log:    conf.Log,
		MenuUC: menuUC,
	})

	// ______________________________________________

	menuxDBRepo := menuxrepo.NewDB(conf.DB)
	menuxUC := menuxuc.New(conf.ServConf, conf.Log, menuxDBRepo)
	menuxweb.Route(w, &menuxweb.Options{
		Log:    conf.Log,
		MenuUC: menuxUC,
	})

	addressxDBRepo := addressxrepo.NewDB(conf.DB)
	addressxUC := addressxuc.New(conf.ServConf, conf.Log, addressxDBRepo)
	addressxweb.Route(w, &addressxweb.Options{
		Log:       conf.Log,
		AddressUC: addressxUC,
	})

	outletxDBRepo := outletxrepo.NewDB(conf.DB)
	outletxUC := outletxuc.New(conf.ServConf, conf.Log, outletxDBRepo)
	outletxweb.Route(w, &outletxweb.Options{
		AddressUC: addressxUC,
		OutletUC:  outletxUC,
		Log:       conf.Log,
	})
}
