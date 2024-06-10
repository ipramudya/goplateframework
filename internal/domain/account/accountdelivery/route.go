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
	con := newController(opts.AccountUC, opts.Log)

	g := web.Echo.Group("/api/v1/account")
	g.POST("/register", con.register)
	g.PUT("/change-password", con.changePassword, web.Mid.Authenticated)
	g.GET("/me", con.me, web.Mid.Authenticated)
}
