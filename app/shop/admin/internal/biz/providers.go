package biz

import "github.com/samber/do/v2"

var ProviderSet = do.Package(
	do.Lazy[*UserUseCase](NewUserUseCase),
	do.Lazy[*CatalogUseCase](NewCatalogUseCase),
)
