package accountdelivery

import (
	"github.com/goplateframework/internal/domain/account/accountuc"
	"github.com/goplateframework/internal/web"
	"github.com/goplateframework/pkg/logger"
)

type Options struct {
	Log       *logger.Log
	AccountUC *accountuc.Usecase
}

func Route(web *web.Web, opts *Options) {
	con := newController(opts.AccountUC)

	g := web.Echo.Group("/api/v1/account")
	g.POST("/register", con.register)
	g.POST("/login", con.login)
	g.PUT("/change-password", con.changePassword, web.Mid.Authenticated)
	g.POST("/logout", con.logout, web.Mid.Authenticated)
}
