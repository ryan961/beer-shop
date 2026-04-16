package server

import (
	"errors"

	"github.com/go-kratos/beer-shop/app/user/service/internal/conf"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/samber/do/v2"

	consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	consulAPI "github.com/hashicorp/consul/api"
)

func NewRegistrar(i do.Injector) (registry.Registrar, error) {
	conf := do.MustInvoke[*conf.Registry](i)
	if conf == nil || conf.Consul == nil {
		return nil, errors.New("registry consul config is required")
	}
	c := consulAPI.DefaultConfig()
	c.Address = conf.Consul.Address
	c.Scheme = conf.Consul.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		return nil, err
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	return r, nil
}
