package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/do/v2"
)

type PaymentRepo interface {
}

type PaymentUseCase struct {
	repo PaymentRepo
	log  *log.Helper
}

func NewPaymentUseCase(i do.Injector) (*PaymentUseCase, error) {
	repo := do.MustInvoke[PaymentRepo](i)
	logger := do.MustInvoke[log.Logger](i)
	return &PaymentUseCase{repo: repo, log: log.NewHelper(log.With(logger, "module", "usecase/payment"))}, nil
}
