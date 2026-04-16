package data

import (
	"errors"

	"github.com/go-kratos/beer-shop/app/shop/admin/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	"context"

	cartv1 "github.com/go-kratos/beer-shop/api/_gen/go/cart/service/v1"
	catalogv1 "github.com/go-kratos/beer-shop/api/_gen/go/catalog/service/v1"
	userv1 "github.com/go-kratos/beer-shop/api/_gen/go/user/service/v1"

	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	consulAPI "github.com/hashicorp/consul/api"
	"github.com/samber/do/v2"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

// Data .
type Data struct {
	log *log.Helper
	uc  userv1.UserClient
	cc  cartv1.CartClient
	bc  catalogv1.CatalogClient
}

// NewData .
func NewData(i do.Injector) (*Data, error) {
	_ = do.MustInvoke[*conf.Data](i)
	logger := do.MustInvoke[log.Logger](i)
	uc := do.MustInvoke[userv1.UserClient](i)
	cc := do.MustInvoke[cartv1.CartClient](i)
	bc := do.MustInvoke[catalogv1.CatalogClient](i)
	l := log.NewHelper(log.With(logger, "module", "data"))
	return &Data{log: l, uc: uc, cc: cc, bc: bc}, nil
}

func NewDiscovery(i do.Injector) (registry.Discovery, error) {
	registryConf := do.MustInvoke[*conf.Registry](i)
	if registryConf == nil || registryConf.Consul == nil {
		return nil, errors.New("registry consul config is required")
	}
	c := consulAPI.DefaultConfig()
	c.Address = registryConf.Consul.Address
	c.Scheme = registryConf.Consul.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		return nil, err
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	return r, nil
}

func NewRegistrar(i do.Injector) (registry.Registrar, error) {
	registryConf := do.MustInvoke[*conf.Registry](i)
	if registryConf == nil || registryConf.Consul == nil {
		return nil, errors.New("registry consul config is required")
	}
	c := consulAPI.DefaultConfig()
	c.Address = registryConf.Consul.Address
	c.Scheme = registryConf.Consul.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		return nil, err
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	return r, nil
}

func NewUserServiceClient(i do.Injector) (userv1.UserClient, error) {
	r := do.MustInvoke[registry.Discovery](i)
	tp := do.MustInvoke[*tracesdk.TracerProvider](i)
	if r == nil {
		return nil, errors.New("service discovery is required")
	}
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///beer.user.service"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			tracing.Client(tracing.WithTracerProvider(tp)),
			recovery.Recovery(),
		),
	)
	if err != nil {
		return nil, err
	}
	c := userv1.NewUserClient(conn)
	return c, nil
}

func NewCartServiceClient(i do.Injector) (cartv1.CartClient, error) {
	r := do.MustInvoke[registry.Discovery](i)
	tp := do.MustInvoke[*tracesdk.TracerProvider](i)
	if r == nil {
		return nil, errors.New("service discovery is required")
	}
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///beer.cart.service"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			tracing.Client(tracing.WithTracerProvider(tp)),
			recovery.Recovery(),
		),
	)
	if err != nil {
		return nil, err
	}
	return cartv1.NewCartClient(conn), nil
}

func NewCatalogServiceClient(i do.Injector) (catalogv1.CatalogClient, error) {
	r := do.MustInvoke[registry.Discovery](i)
	tp := do.MustInvoke[*tracesdk.TracerProvider](i)
	if r == nil {
		return nil, errors.New("service discovery is required")
	}
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///beer.catalog.service"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			tracing.Client(tracing.WithTracerProvider(tp)),
			recovery.Recovery(),
		),
	)
	if err != nil {
		return nil, err
	}
	return catalogv1.NewCatalogClient(conn), nil
}
