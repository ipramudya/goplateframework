package outletweb

import (
	"github.com/goplateframework/internal/web"
	"github.com/goplateframework/pkg/logger"
)

type Options struct {
	AddressUC iAddressUsecase
	OutletUC  iOutletUsecase
	Log       *logger.Log
}

func Route(web *web.Web, opts *Options) {
	con := newController(opts.AddressUC, opts.OutletUC, opts.Log)

	g := web.Echo.Group("/api/v1/outlet", web.Mid.Authenticated)
	g.GET("", con.getAll)
	g.GET("/:id", con.getOne)
	g.POST("", con.create)
	g.PUT("/:id", con.update)
}
