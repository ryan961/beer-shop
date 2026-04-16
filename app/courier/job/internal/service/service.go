package service

import (
	v1 "github.com/go-kratos/beer-shop/api/_gen/go/courier/job/v1"
	"github.com/go-kratos/beer-shop/app/courier/job/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/do/v2"
)

type CourierService struct {
	v1.UnimplementedCourierServer

	oc  *biz.CourierUseCase
	log *log.Helper
}

func NewCourierService(i do.Injector) (*CourierService, error) {
	oc := do.MustInvoke[*biz.CourierUseCase](i)
	logger := do.MustInvoke[log.Logger](i)
	return &CourierService{
		oc:  oc,
		log: log.NewHelper(log.With(logger, "module", "service/courier")),
	}, nil
}
