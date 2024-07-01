package authweb

import (
	"github.com/goplateframework/internal/web"
	"github.com/goplateframework/pkg/logger"
)

type Options struct {
	Log    *logger.Log
	AuthUC iUsecase
}

func Route(web *web.Web, opts *Options) {
	con := newController(opts.AuthUC, opts.Log)

	g := web.Echo.Group("/api/v1/new-auth")
	g.POST("/login", con.login)
	g.POST("/logout", con.logout, web.Mid.RefreshAuth, web.Mid.Authenticated)
	g.POST("/refresh", con.refreshToken, web.Mid.RefreshAuth, web.Mid.Authenticated)
}
