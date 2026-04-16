package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/do/v2"
)

type ShipOrder struct {
	Id     int64
	UserId int64
}

type CourierRepo interface {
}

type CourierUseCase struct {
	repo CourierRepo
	log  *log.Helper
}

func NewCourierUseCase(i do.Injector) (*CourierUseCase, error) {
	repo := do.MustInvoke[CourierRepo](i)
	logger := do.MustInvoke[log.Logger](i)
	return &CourierUseCase{repo: repo, log: log.NewHelper(log.With(logger, "module", "usecase/courier"))}, nil
}
