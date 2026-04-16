package server

import (
	"testing"

	"github.com/go-kratos/beer-shop/app/shop/interface/internal/conf"
	"github.com/samber/do/v2"
)

func TestNewRegistrarRequiresConfig(t *testing.T) {
	t.Parallel()

	injector := do.New(func(i do.Injector) {
		do.ProvideValue(i, (*conf.Registry)(nil))
	})

	registrar, err := NewRegistrar(injector)
	if err == nil {
		t.Fatal("expected error for nil registry config")
	}
	if registrar != nil {
		t.Fatal("expected nil registrar for nil registry config")
	}
}

func TestNewRegistrarBuildsConsulClient(t *testing.T) {
	t.Parallel()

	injector := do.New(func(i do.Injector) {
		do.ProvideValue(i, &conf.Registry{
			Consul: &conf.Registry_Consul{
				Address: "127.0.0.1:8500",
				Scheme:  "http",
			},
		})
	})

	registrar, err := NewRegistrar(injector)
	if err != nil {
		t.Fatalf("expected registrar to be created: %v", err)
	}
	if registrar == nil {
		t.Fatal("expected non-nil registrar")
	}
}
