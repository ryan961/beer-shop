package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/do/v2"

	"github.com/go-kratos/beer-shop/app/courier/job/internal/biz"
)

var _ biz.CourierRepo = (*courierRepo)(nil)

type courierRepo struct {
	data *Data
	log  *log.Helper
}

type ShippingEntry struct {
	OrderId string `json:"order_id"`
}

func NewCourierRepo(i do.Injector) (biz.CourierRepo, error) {
	data := do.MustInvoke[*Data](i)
	logger := do.MustInvoke[log.Logger](i)
	return &courierRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "data/courier")),
	}, nil
}
