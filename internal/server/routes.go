package server

import (
	"github.com/goplateframework/internal/domain/account/accountrepo"
	"github.com/goplateframework/internal/domain/account/accountuc"
	"github.com/goplateframework/internal/domain/account/accountweb"
	"github.com/goplateframework/internal/domain/address/addressrepo"
	"github.com/goplateframework/internal/domain/address/addressuc"
	"github.com/goplateframework/internal/domain/address/addressweb"
	"github.com/goplateframework/internal/domain/auth/authrepo"
	"github.com/goplateframework/internal/domain/auth/authuc"
	"github.com/goplateframework/internal/domain/auth/authweb"
	"github.com/goplateframework/internal/domain/menu/menurepo"
	"github.com/goplateframework/internal/domain/menu/menuuc"
	"github.com/goplateframework/internal/domain/menu/menuweb"
	"github.com/goplateframework/internal/domain/outlet/outletrepo"
	"github.com/goplateframework/internal/domain/outlet/outletuc"
	"github.com/goplateframework/internal/domain/outlet/outletweb"
	"github.com/goplateframework/internal/web"
)

func routes(w *web.Web, conf *Config) {
	authCacheRepo := authrepo.NewCache(conf.RDB)
	accountDBRepo := accountrepo.NewDB(conf.DB)
	accountCacheRepo := accountrepo.NewCache(conf.RDB)

	accountUC := accountuc.New(conf.ServConf, conf.Log, accountDBRepo, accountCacheRepo)
	authUC := authuc.New(conf.ServConf, conf.Log, authCacheRepo, accountDBRepo)

	accountweb.Route(w, &accountweb.Options{
		Log:       conf.Log,
		AccountUC: accountUC,
	})

	authweb.Route(w, &authweb.Options{
		Log:    conf.Log,
		AuthUC: authUC,
	})

	// ______________________________________________

	menuDBRepo := menurepo.NewDB(conf.DB)
	menuUC := menuuc.New(conf.ServConf, conf.Log, menuDBRepo)
	menuweb.Route(w, &menuweb.Options{
		Log:    conf.Log,
		MenuUC: menuUC,
	})

	addressDBRepo := addressrepo.NewDB(conf.DB)
	addressUC := addressuc.New(conf.ServConf, conf.Log, addressDBRepo)
	addressweb.Route(w, &addressweb.Options{
		Log:       conf.Log,
		AddressUC: addressUC,
	})

	outletDBRepo := outletrepo.NewDB(conf.DB)
	outletUC := outletuc.New(conf.ServConf, conf.Log, outletDBRepo)
	outletweb.Route(w, &outletweb.Options{
		AddressUC: addressUC,
		OutletUC:  outletUC,
		Log:       conf.Log,
	})
}
