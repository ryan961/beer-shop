package server

import (
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	httptransport "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/samber/do/v2"
)

var ProviderSet = do.Package(
	do.Lazy[*httptransport.Server](NewHTTPServer),
	do.Lazy[*grpc.Server](NewGRPCServer),
	do.Lazy[registry.Registrar](NewRegistrar),
)
