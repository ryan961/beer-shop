package service

import (
	v1 "github.com/go-kratos/beer-shop/api/_gen/go/shipping/service/v1"
	"github.com/go-kratos/beer-shop/app/shipping/service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/do/v2"
)

type ShippingService struct {
	v1.UnimplementedShippingServer

	oc  *biz.ShippingUseCase
	log *log.Helper
}

func NewShippingService(i do.Injector) (*ShippingService, error) {
	oc := do.MustInvoke[*biz.ShippingUseCase](i)
	logger := do.MustInvoke[log.Logger](i)
	return &ShippingService{
		oc:  oc,
		log: log.NewHelper(log.With(logger, "module", "service/shipping")),
	}, nil
}
