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
	"github.com/goplateframework/internal/domain/menutoping/menutopingrepo"
	"github.com/goplateframework/internal/domain/menutoping/menutopinguc"
	"github.com/goplateframework/internal/domain/menutoping/menutopingweb"
	"github.com/goplateframework/internal/domain/outlet/outletrepo"
	"github.com/goplateframework/internal/domain/outlet/outletuc"
	"github.com/goplateframework/internal/domain/outlet/outletweb"
	"github.com/goplateframework/internal/web"
)

func router(w *web.Web, conf *Config) {
	accountDBRepo := accountrepo.NewDB(conf.DB)
	accountCacheRepo := accountrepo.NewCache(conf.RDB)
	accountUC := accountuc.New(conf.ServConf, conf.Log, accountDBRepo, accountCacheRepo)
	accountweb.Route(w, &accountweb.Options{
		Log:       conf.Log,
		AccountUC: accountUC,
	})

	authCacheRepo := authrepo.NewCache(conf.RDB)
	authUC := authuc.New(conf.ServConf, conf.Log, authCacheRepo, accountDBRepo)
	authweb.Route(w, &authweb.Options{
		Log:    conf.Log,
		AuthUC: authUC,
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

	menuDBRepo := menurepo.NewDB(conf.DB)
	menuUC := menuuc.New(conf.ServConf, conf.Log, conf.Worker, menuDBRepo)
	menuweb.Route(w, &menuweb.Options{
		Log:    conf.Log,
		MenuUC: menuUC,
	})

	menuTopingDBRepo := menutopingrepo.NewDB(conf.DB)
	menuTopingUC := menutopinguc.New(conf.ServConf, conf.Log, menuTopingDBRepo)
	menutopingweb.Route(w, &menutopingweb.Options{
		Log:          conf.Log,
		MenuTopingUC: menuTopingUC,
	})
}
