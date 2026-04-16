package data

import (
	"github.com/go-kratos/beer-shop/app/payment/service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/do/v2"
)

var _ biz.PaymentRepo = (*paymentRepo)(nil)

type paymentRepo struct {
	data *Data
	log  *log.Helper
}

func NewPaymentRepo(i do.Injector) (biz.PaymentRepo, error) {
	data := do.MustInvoke[*Data](i)
	logger := do.MustInvoke[log.Logger](i)
	return &paymentRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/payment")),
	}, nil
}
