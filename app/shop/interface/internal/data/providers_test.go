package data

import (
	"testing"

	"github.com/go-kratos/beer-shop/app/shop/interface/internal/conf"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/samber/do/v2"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func TestNewDiscoveryRequiresConfig(t *testing.T) {
	t.Parallel()

	injector := do.New(func(i do.Injector) {
		do.ProvideValue(i, (*conf.Registry)(nil))
	})

	discovery, err := NewDiscovery(injector)
	if err == nil {
		t.Fatal("expected error for nil registry config")
	}
	if discovery != nil {
		t.Fatal("expected nil discovery for nil registry config")
	}
}

func TestNewUserServiceClientRequiresAuth(t *testing.T) {
	t.Parallel()

	injector := do.New(func(i do.Injector) {
		do.ProvideValue(i, (*conf.Auth)(nil))

		var discovery registry.Discovery
		do.ProvideValue(i, discovery)
		do.ProvideValue(i, (*tracesdk.TracerProvider)(nil))
	})

	client, err := NewUserServiceClient(injector)
	if err == nil {
		t.Fatal("expected error for nil auth config")
	}
	if client != nil {
		t.Fatal("expected nil client for nil auth config")
	}
}

func TestNewDiscoveryBuildsConsulClient(t *testing.T) {
	t.Parallel()

	injector := do.New(func(i do.Injector) {
		do.ProvideValue(i, &conf.Registry{
			Consul: &conf.Registry_Consul{
				Address: "127.0.0.1:8500",
				Scheme:  "http",
			},
		})
	})

	discovery, err := NewDiscovery(injector)
	if err != nil {
		t.Fatalf("expected discovery to be created: %v", err)
	}
	if discovery == nil {
		t.Fatal("expected non-nil discovery")
	}
}
