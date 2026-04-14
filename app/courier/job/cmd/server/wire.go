package main

import (
	orderv1 "github.com/go-kratos/beer-shop/api/_gen/go/order/service/v1"
	"github.com/go-kratos/beer-shop/app/courier/job/internal/biz"
	"github.com/go-kratos/beer-shop/app/courier/job/internal/conf"
	"github.com/go-kratos/beer-shop/app/courier/job/internal/data"
	"github.com/go-kratos/beer-shop/app/courier/job/internal/server"
	"github.com/go-kratos/beer-shop/app/courier/job/internal/service"

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
			do.Provide(i, func(i do.Injector) (registry.Discovery, error) {
				return data.NewDiscovery(registryConfig), nil
			})
			do.Provide(i, func(i do.Injector) (orderv1.OrderClient, error) {
				return data.NewOrderServiceClient(do.MustInvoke[registry.Discovery](i), tracerProvider), nil
			})
			do.Provide(i, func(i do.Injector) (*data.Data, error) {
				return data.NewData(data.NewKafkaConsumer(confData), logger, do.MustInvoke[orderv1.OrderClient](i))
			})
			do.Provide(i, func(i do.Injector) (biz.CourierRepo, error) {
				return data.NewCourierRepo(do.MustInvoke[*data.Data](i), logger), nil
			})
			do.Provide(i, func(i do.Injector) (*biz.CourierUseCase, error) {
				return biz.NewCourierUseCase(do.MustInvoke[biz.CourierRepo](i), logger), nil
			})
			do.Provide(i, func(i do.Injector) (*service.CourierService, error) {
				return service.NewCourierService(do.MustInvoke[*biz.CourierUseCase](i), logger), nil
			})
			do.Provide(i, func(i do.Injector) (*grpc.Server, error) {
				return server.NewGRPCServer(confServer, logger, tracerProvider, do.MustInvoke[*service.CourierService](i)), nil
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
