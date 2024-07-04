package menutopingweb

import (
	"github.com/goplateframework/internal/web"
	"github.com/goplateframework/pkg/logger"
)

type Options struct {
	Log          *logger.Log
	MenuTopingUC iUsecase
}

func Route(web *web.Web, opts *Options) {
	con := newController(opts.MenuTopingUC, opts.Log)

	g := web.Echo.Group("/api/v1/menu-topings", web.Mid.Authenticated)
	g.POST("", con.create)
	g.GET("", con.getAll)
	g.GET("/:id", con.getOne)
	g.PUT("/:id", con.update)
	g.DELETE("/:id", con.delete)
}
