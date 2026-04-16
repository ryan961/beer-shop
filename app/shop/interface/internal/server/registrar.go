package server

import (
	"errors"

	"github.com/go-kratos/beer-shop/app/shop/interface/internal/conf"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/registry"
	consulAPI "github.com/hashicorp/consul/api"
	"github.com/samber/do/v2"
)

func NewRegistrar(i do.Injector) (registry.Registrar, error) {
	registryConfig := do.MustInvoke[*conf.Registry](i)
	if registryConfig == nil || registryConfig.Consul == nil {
		return nil, errors.New("registry consul config is required")
	}

	c := consulAPI.DefaultConfig()
	c.Address = registryConfig.Consul.Address
	c.Scheme = registryConfig.Consul.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		return nil, err
	}

	return consul.New(cli, consul.WithHealthCheck(false)), nil
}
