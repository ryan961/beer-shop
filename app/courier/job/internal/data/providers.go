package data

import (
	orderv1 "github.com/go-kratos/beer-shop/api/_gen/go/order/service/v1"
	"github.com/go-kratos/beer-shop/app/courier/job/internal/biz"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/samber/do/v2"
)

var ProviderSet = do.Package(
	do.Lazy[registry.Discovery](NewDiscovery),
	do.Lazy[orderv1.OrderClient](NewOrderServiceClient),
	do.Lazy[*Data](NewData),
	do.Lazy[biz.CourierRepo](NewCourierRepo),
)
