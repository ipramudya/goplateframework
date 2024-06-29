package menuxweb

import (
	"github.com/goplateframework/internal/web"
	"github.com/goplateframework/pkg/logger"
)

type Options struct {
	Log    *logger.Log
	MenuUC Usecase
}

func Route(web *web.Web, opts *Options) {
	con := newController(opts.MenuUC, opts.Log)

	g := web.Echo.Group("/api/v1/new-menu", web.Mid.Authenticated)
	g.POST("", con.create)
	g.GET("", con.getAll)
	g.PUT("/:id", con.update)
}
