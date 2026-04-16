package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/go-kratos/beer-shop/app/shop/interface/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/samber/do/v2"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.4.0"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "beer.shop.interface"
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(i do.Injector) (*kratos.App, error) {
	logger := do.MustInvoke[log.Logger](i)
	hs := do.MustInvoke[*http.Server](i)
	gs := do.MustInvoke[*grpc.Server](i)
	rr := do.MustInvoke[registry.Registrar](i)

	return kratos.New(
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			gs,
		),
		kratos.Registrar(rr),
	), nil
}

func main() {
	logger := log.With(log.NewStdLogger(os.Stdout),
		"service.name", Name,
		"service.version", Version,
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
	)
	if err := run(logger); err != nil {
		log.NewHelper(logger).Error(err)
		os.Exit(1)
	}
}

func run(logger log.Logger) error {
	flag.Parse()

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer func() {
		_ = c.Close()
	}()
	if err := c.Load(); err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		return fmt.Errorf("scan bootstrap config: %w", err)
	}

	var rc conf.Registry
	if err := c.Scan(&rc); err != nil {
		return fmt.Errorf("scan registry config: %w", err)
	}
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(bc.Trace.Endpoint)))
	if err != nil {
		return fmt.Errorf("create jaeger exporter: %w", err)
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(Name),
		)),
	)
	defer func() {
		_ = tp.Shutdown(context.Background())
	}()

	app, cleanup, err := initApp(bc.Server, &rc, bc.Data, bc.Auth, logger, tp)
	if err != nil {
		return fmt.Errorf("init app: %w", err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		return fmt.Errorf("run app: %w", err)
	}

	return nil
}
