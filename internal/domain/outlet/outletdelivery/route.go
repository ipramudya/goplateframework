package outletdelivery

import (
	"github.com/goplateframework/internal/domain/outlet/outletuc"
	"github.com/goplateframework/internal/web"
	"github.com/goplateframework/pkg/logger"
)

type Options struct {
	Log      *logger.Log
	OutletUC *outletuc.Usecase
}

func Route(web *web.Web, opts *Options) {
	con := newController(opts.OutletUC, opts.Log)

	g := web.Echo.Group("/api/v1/outlet", web.Mid.Authenticated)
	g.POST("", con.create)
	g.GET("/:id", con.getOne)
}
