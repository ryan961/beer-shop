package biz

import "github.com/samber/do/v2"

var ProviderSet = do.Package(
	do.Lazy[*PaymentUseCase](NewPaymentUseCase),
)
