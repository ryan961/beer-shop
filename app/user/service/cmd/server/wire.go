package main

import (
	"github.com/go-kratos/beer-shop/app/user/service/internal/biz"
	"github.com/go-kratos/beer-shop/app/user/service/internal/conf"
	"github.com/go-kratos/beer-shop/app/user/service/internal/data"
	"github.com/go-kratos/beer-shop/app/user/service/internal/server"
	"github.com/go-kratos/beer-shop/app/user/service/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/samber/do/v2"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

// initApp init kratos application.
func initApp(confServer *conf.Server, registryConfig *conf.Registry, confData *conf.Data, auth *conf.Auth, logger log.Logger, tracerProvider *tracesdk.TracerProvider) (*kratos.App, func(), error) {
	injector := do.New(
		func(i do.Injector) {
			do.Provide(i, func(i do.Injector) (*data.Data, error) {
				return data.NewData(data.NewEntClient(confData, logger), data.NewRedisCmd(confData, logger), logger)
			})
			do.Provide(i, func(i do.Injector) (biz.UserRepo, error) {
				return data.NewUserRepo(do.MustInvoke[*data.Data](i), logger), nil
			})
			do.Provide(i, func(i do.Injector) (biz.CardRepo, error) {
				return data.NewCardRepo(do.MustInvoke[*data.Data](i), logger), nil
			})
			do.Provide(i, func(i do.Injector) (biz.AddressRepo, error) {
				return data.NewAddressRepo(do.MustInvoke[*data.Data](i), logger), nil
			})
			do.Provide(i, func(i do.Injector) (*biz.UserUseCase, error) {
				return biz.NewUserUseCase(do.MustInvoke[biz.UserRepo](i), logger), nil
			})
			do.Provide(i, func(i do.Injector) (*biz.CardUseCase, error) {
				return biz.NewCardUseCase(do.MustInvoke[biz.CardRepo](i), logger), nil
			})
			do.Provide(i, func(i do.Injector) (*biz.AddressUseCase, error) {
				return biz.NewAddressUseCase(do.MustInvoke[biz.AddressRepo](i), logger), nil
			})
			do.Provide(i, func(i do.Injector) (*service.UserService, error) {
				return service.NewUserService(
					do.MustInvoke[*biz.UserUseCase](i),
					do.MustInvoke[*biz.CardUseCase](i),
					do.MustInvoke[*biz.AddressUseCase](i),
					logger,
				), nil
			})
			do.Provide(i, func(i do.Injector) (*grpc.Server, error) {
				return server.NewGRPCServer(confServer, auth, logger, tracerProvider, do.MustInvoke[*service.UserService](i)), nil
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
