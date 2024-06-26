package menudelivery

import (
	"github.com/goplateframework/internal/domain/menu/menuuc"
	"github.com/goplateframework/internal/web"
	"github.com/goplateframework/pkg/logger"
)

type Options struct {
	Log    *logger.Log
	MenuUC *menuuc.Usecase
}

func Route(web *web.Web, opts *Options) {
	con := newController(opts.MenuUC, opts.Log)

	g := web.Echo.Group("/api/v1/menu", web.Mid.Authenticated)
	g.POST("", con.create)
	g.PUT("/:id", con.update)
	g.GET("/:outletId", con.getAllByOutletID)
}
