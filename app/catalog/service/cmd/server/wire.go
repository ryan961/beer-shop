package main

import (
	"github.com/go-kratos/beer-shop/app/catalog/service/internal/biz"
	"github.com/go-kratos/beer-shop/app/catalog/service/internal/conf"
	"github.com/go-kratos/beer-shop/app/catalog/service/internal/data"
	"github.com/go-kratos/beer-shop/app/catalog/service/internal/server"
	"github.com/go-kratos/beer-shop/app/catalog/service/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/do/v2"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

// initApp init kratos application.
func initApp(confServer *conf.Server, registryConfig *conf.Registry, confData *conf.Data, logger log.Logger, tracerProvider *tracesdk.TracerProvider) (*kratos.App, func(), error) {
	injector := do.New(
		do.Eager[*conf.Server](confServer),
		do.Eager[*conf.Registry](registryConfig),
		do.Eager[*conf.Data](confData),
		do.Eager[log.Logger](logger),
		do.Eager[*tracesdk.TracerProvider](tracerProvider),
		data.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		server.ProviderSet,
		do.Package(do.Lazy[*kratos.App](newApp)),
	)

	app, err := do.Invoke[*kratos.App](injector)
	if err != nil {
		return nil, nil, err
	}

	return app, func() {
		_ = injector.Shutdown()
	}, nil
}
