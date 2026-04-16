package service

import (
	v1 "github.com/go-kratos/beer-shop/api/_gen/go/payment/service/v1"
	"github.com/go-kratos/beer-shop/app/payment/service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/do/v2"
)

type PaymentService struct {
	v1.UnimplementedPaymentServer

	pc  *biz.PaymentUseCase
	log *log.Helper
}

func NewPaymentService(i do.Injector) (*PaymentService, error) {
	pc := do.MustInvoke[*biz.PaymentUseCase](i)
	logger := do.MustInvoke[log.Logger](i)
	return &PaymentService{
		pc:  pc,
		log: log.NewHelper(log.With(logger, "module", "service/payment")),
	}, nil
}
