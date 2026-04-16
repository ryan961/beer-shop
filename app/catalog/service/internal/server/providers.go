package server

import (
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/samber/do/v2"
)

var ProviderSet = do.Package(
	do.Lazy[*grpc.Server](NewGRPCServer),
	do.Lazy[registry.Registrar](NewRegistrar),
)
