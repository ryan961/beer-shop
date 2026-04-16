package server

import (
	v1 "github.com/go-kratos/beer-shop/api/_gen/go/shop/admin/v1"
	"github.com/go-kratos/beer-shop/app/shop/admin/internal/conf"
	"github.com/go-kratos/beer-shop/app/shop/admin/internal/service"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/gorilla/handlers"
	"github.com/samber/do/v2"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(i do.Injector) (*http.Server, error) {
	c := do.MustInvoke[*conf.Server](i)
	logger := do.MustInvoke[log.Logger](i)
	tp := do.MustInvoke[*tracesdk.TracerProvider](i)
	s := do.MustInvoke[*service.ShopAdmin](i)
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(
				tracing.WithTracerProvider(tp)),
			logging.Server(logger),
		),
		http.Filter(handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}),
		)),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterShopAdminHTTPServer(srv, s)
	return srv, nil
}
