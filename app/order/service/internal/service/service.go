package service

import (
	v1 "github.com/go-kratos/beer-shop/api/_gen/go/order/service/v1"
	"github.com/go-kratos/beer-shop/app/order/service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type OrderService struct {
	v1.UnimplementedOrderServer

	oc  *biz.OrderUseCase
	log *log.Helper
}

func NewOrderService(oc *biz.OrderUseCase, logger log.Logger) *OrderService {
	return &OrderService{
		oc:  oc,
		log: log.NewHelper(log.With(logger, "module", "service/order"))}
}
