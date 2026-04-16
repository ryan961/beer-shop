package service

import (
	v1 "github.com/go-kratos/beer-shop/api/_gen/go/order/service/v1"
	"github.com/go-kratos/beer-shop/app/order/service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/do/v2"
)

type OrderService struct {
	v1.UnimplementedOrderServer

	oc  *biz.OrderUseCase
	log *log.Helper
}

func NewOrderService(i do.Injector) (*OrderService, error) {
	oc := do.MustInvoke[*biz.OrderUseCase](i)
	logger := do.MustInvoke[log.Logger](i)
	return &OrderService{
		oc:  oc,
		log: log.NewHelper(log.With(logger, "module", "service/order")),
	}, nil
}
