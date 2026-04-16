package data

import (
	"github.com/go-kratos/beer-shop/app/user/service/internal/biz"
	"github.com/samber/do/v2"
)

var ProviderSet = do.Package(
	do.Lazy[*Data](NewData),
	do.Lazy[biz.UserRepo](NewUserRepo),
	do.Lazy[biz.CardRepo](NewCardRepo),
	do.Lazy[biz.AddressRepo](NewAddressRepo),
)
