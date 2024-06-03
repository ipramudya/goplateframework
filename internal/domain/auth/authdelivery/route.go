package authdelivery

import (
	"github.com/goplateframework/internal/domain/auth/authuc"
	"github.com/goplateframework/internal/web"
	"github.com/goplateframework/pkg/logger"
)

type Options struct {
	Log    *logger.Log
	AuthUC *authuc.Usecase
}

func Route(web *web.Web, opts *Options) {
	con := newController(opts.AuthUC)

	g := web.Echo.Group("/api/v1/auth")
	g.POST("/login", con.login)
	g.POST("/logout", con.logout, web.Mid.Authenticated)
}
