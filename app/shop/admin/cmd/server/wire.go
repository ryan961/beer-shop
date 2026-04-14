package main

import (
	cartv1 "github.com/go-kratos/beer-shop/api/_gen/go/cart/service/v1"
	catalogv1 "github.com/go-kratos/beer-shop/api/_gen/go/catalog/service/v1"
	userv1 "github.com/go-kratos/beer-shop/api/_gen/go/user/service/v1"
	"github.com/go-kratos/beer-shop/app/shop/admin/internal/biz"
	"github.com/go-kratos/beer-shop/app/shop/admin/internal/conf"
	"github.com/go-kratos/beer-shop/app/shop/admin/internal/data"
	"github.com/go-kratos/beer-shop/app/shop/admin/internal/server"
	"github.com/go-kratos/beer-shop/app/shop/admin/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	httptransport "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/samber/do/v2"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

// initApp init kratos application.
func initApp(confServer *conf.Server, registryConfig *conf.Registry, confData *conf.Data, logger log.Logger, tracerProvider *tracesdk.TracerProvider) (*kratos.App, func(), error) {
	injector := do.New(
		func(i do.Injector) {
			do.Provide(i, func(i do.Injector) (registry.Discovery, error) {
				return data.NewDiscovery(registryConfig), nil
			})
			do.Provide(i, func(i do.Injector) (userv1.UserClient, error) {
				return data.NewUserServiceClient(do.MustInvoke[registry.Discovery](i), tracerProvider), nil
			})
			do.Provide(i, func(i do.Injector) (cartv1.CartClient, error) {
				return data.NewCartServiceClient(do.MustInvoke[registry.Discovery](i), tracerProvider), nil
			})
			do.Provide(i, func(i do.Injector) (catalogv1.CatalogClient, error) {
				return data.NewCatalogServiceClient(do.MustInvoke[registry.Discovery](i), tracerProvider), nil
			})
			do.Provide(i, func(i do.Injector) (*data.Data, error) {
				return data.NewData(
					confData,
					logger,
					do.MustInvoke[userv1.UserClient](i),
					do.MustInvoke[cartv1.CartClient](i),
					do.MustInvoke[catalogv1.CatalogClient](i),
				)
			})
			do.Provide(i, func(i do.Injector) (biz.UserRepo, error) {
				return data.NewUserRepo(do.MustInvoke[*data.Data](i), logger), nil
			})
			do.Provide(i, func(i do.Injector) (biz.CatalogRepo, error) {
				return data.NewCatalogRepo(do.MustInvoke[*data.Data](i), logger), nil
			})
			do.Provide(i, func(i do.Injector) (*biz.UserUseCase, error) {
				return biz.NewUserUseCase(do.MustInvoke[biz.UserRepo](i), logger, do.MustInvoke[userv1.UserClient](i)), nil
			})
			do.Provide(i, func(i do.Injector) (*biz.CatalogUseCase, error) {
				return biz.NewCatalogUseCase(do.MustInvoke[biz.CatalogRepo](i), logger), nil
			})
			do.Provide(i, func(i do.Injector) (*service.ShopAdmin, error) {
				return service.NewShopAdmin(do.MustInvoke[*biz.UserUseCase](i), do.MustInvoke[*biz.CatalogUseCase](i), logger), nil
			})
			do.Provide(i, func(i do.Injector) (*httptransport.Server, error) {
				return server.NewHTTPServer(confServer, logger, tracerProvider, do.MustInvoke[*service.ShopAdmin](i)), nil
			})
			do.Provide(i, func(i do.Injector) (*grpc.Server, error) {
				return server.NewGRPCServer(confServer, logger, tracerProvider, do.MustInvoke[*service.ShopAdmin](i)), nil
			})
			do.Provide(i, func(i do.Injector) (registry.Registrar, error) {
				return data.NewRegistrar(registryConfig), nil
			})
			do.Provide(i, func(i do.Injector) (*kratos.App, error) {
				return newApp(
					logger,
					do.MustInvoke[*httptransport.Server](i),
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
