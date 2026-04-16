package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/do/v2"
)

type ShipOrder struct {
	Id     int64
	UserId int64
}

type ShippingRepo interface {
	ShipOrder(ctx context.Context, o *ShipOrder) error
}

type ShippingUseCase struct {
	repo ShippingRepo
	log  *log.Helper
}

func NewShippingUseCase(i do.Injector) (*ShippingUseCase, error) {
	repo := do.MustInvoke[ShippingRepo](i)
	logger := do.MustInvoke[log.Logger](i)
	return &ShippingUseCase{repo: repo, log: log.NewHelper(log.With(logger, "module", "usecase/shipping"))}, nil
}

func (uc *ShippingUseCase) ShipOrder(ctx context.Context, o *ShipOrder) (err error) {
	return uc.repo.ShipOrder(ctx, o)
}
