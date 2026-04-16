package data

import (
	"github.com/go-kratos/beer-shop/app/order/service/internal/biz"
	"github.com/samber/do/v2"
)

var ProviderSet = do.Package(
	do.Lazy[*Data](NewData),
	do.Lazy[biz.OrderRepo](NewOrderRepo),
)
