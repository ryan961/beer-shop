package data

import (
	cartv1 "github.com/go-kratos/beer-shop/api/_gen/go/cart/service/v1"
	catalogv1 "github.com/go-kratos/beer-shop/api/_gen/go/catalog/service/v1"
	userv1 "github.com/go-kratos/beer-shop/api/_gen/go/user/service/v1"
	"github.com/go-kratos/beer-shop/app/shop/admin/internal/biz"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/samber/do/v2"
)

var ProviderSet = do.Package(
	do.Lazy[registry.Discovery](NewDiscovery),
	do.Lazy[userv1.UserClient](NewUserServiceClient),
	do.Lazy[cartv1.CartClient](NewCartServiceClient),
	do.Lazy[catalogv1.CatalogClient](NewCatalogServiceClient),
	do.Lazy[*Data](NewData),
	do.Lazy[biz.UserRepo](NewUserRepo),
	do.Lazy[biz.CatalogRepo](NewCatalogRepo),
	do.Lazy[registry.Registrar](NewRegistrar),
)
