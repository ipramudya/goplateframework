package addressweb

import (
	"github.com/goplateframework/internal/web"
	"github.com/goplateframework/pkg/logger"
)

type Options struct {
	Log       *logger.Log
	AddressUC iUsecase
}

func Route(web *web.Web, opts *Options) {
	con := newController(opts.AddressUC, opts.Log)

	g := web.Echo.Group("/api/v1/address", web.Mid.Authenticated)
	g.PUT("/:id", con.update)
}
