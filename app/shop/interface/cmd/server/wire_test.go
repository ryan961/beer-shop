package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-kratos/beer-shop/app/shop/interface/internal/conf"
	"github.com/go-kratos/kratos/v2/log"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func TestInitAppBuildsInjectorFromProviderSets(t *testing.T) {
	t.Parallel()

	consulServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("[]"))
	}))
	t.Cleanup(consulServer.Close)

	consulURL, err := url.Parse(consulServer.URL)
	if err != nil {
		t.Fatalf("parse consul server URL: %v", err)
	}

	app, cleanup, err := initApp(
		&conf.Server{
			Http: &conf.Server_HTTP{Addr: "127.0.0.1:0"},
			Grpc: &conf.Server_GRPC{Addr: "127.0.0.1:0"},
		},
		&conf.Registry{
			Consul: &conf.Registry_Consul{
				Address: consulURL.Host,
				Scheme:  consulURL.Scheme,
			},
		},
		&conf.Data{},
		&conf.Auth{
			ServiceKey: "service-key",
			ApiKey:     "api-key",
		},
		log.NewStdLogger(io.Discard),
		tracesdk.NewTracerProvider(),
	)
	if err != nil {
		t.Fatalf("expected app to build: %v", err)
	}
	if app == nil {
		t.Fatal("expected non-nil app")
	}
	if cleanup == nil {
		t.Fatal("expected non-nil cleanup")
	}

	cleanup()
}
