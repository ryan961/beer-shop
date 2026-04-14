package main

import (
	"github.com/go-kratos/beer-shop/app/catalog/service/internal/biz"
	"github.com/go-kratos/beer-shop/app/catalog/service/internal/conf"
	"github.com/go-kratos/beer-shop/app/catalog/service/internal/data"
	"github.com/go-kratos/beer-shop/app/catalog/service/internal/server"
	"github.com/go-kratos/beer-shop/app/catalog/service/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/samber/do/v2"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

// initApp init kratos application.
func initApp(confServer *conf.Server, registryConfig *conf.Registry, confData *conf.Data, logger log.Logger, tracerProvider *tracesdk.TracerProvider) (*kratos.App, func(), error) {
	injector := do.New(
		func(i do.Injector) {
			do.Provide(i, func(i do.Injector) (*data.Data, error) {
				return data.NewData(data.NewEntClient(confData, logger), logger)
			})
			do.Provide(i, func(i do.Injector) (biz.BeerRepo, error) {
				return data.NewBeerRepo(do.MustInvoke[*data.Data](i), logger), nil
			})
			do.Provide(i, func(i do.Injector) (*biz.BeerUseCase, error) {
				return biz.NewBeerUseCase(do.MustInvoke[biz.BeerRepo](i), logger), nil
			})
			do.Provide(i, func(i do.Injector) (*service.CatalogService, error) {
				return service.NewCatalogService(do.MustInvoke[*biz.BeerUseCase](i), logger), nil
			})
			do.Provide(i, func(i do.Injector) (*grpc.Server, error) {
				return server.NewGRPCServer(confServer, logger, tracerProvider, do.MustInvoke[*service.CatalogService](i)), nil
			})
			do.Provide(i, func(i do.Injector) (registry.Registrar, error) {
				return server.NewRegistrar(registryConfig), nil
			})
			do.Provide(i, func(i do.Injector) (*kratos.App, error) {
				return newApp(
					logger,
					do.MustInvoke[*grpc.Server](i),
					do.MustInvoke[registry.Registrar](i),
				), nil
			})
		},
	)

	app, err := do.Invoke[*kratos.App](injector)
	if err != nil {
		return nil, nil, err
	}

	return app, func() {
		injector.Shutdown()
	}, nil
}
