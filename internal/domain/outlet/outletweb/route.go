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

	g := web.Echo.Group("/api/v1/new-outlet", web.Mid.Authenticated)
	g.POST("", con.create)
	g.GET("/:id", con.getOne)
	g.PUT("/:id", con.update)
}
