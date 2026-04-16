package data

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/do/v2"

	"github.com/go-kratos/beer-shop/app/shipping/service/internal/biz"
)

var _ biz.ShippingRepo = (*shippingRepo)(nil)

type shippingRepo struct {
	data *Data
	log  *log.Helper
}

func NewShippingRepo(i do.Injector) (biz.ShippingRepo, error) {
	data := do.MustInvoke[*Data](i)
	logger := do.MustInvoke[log.Logger](i)
	return &shippingRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/shipping")),
	}, nil
}

type ShippingEntry struct {
	OrderId string `json:"order_id"`
}

func (uc *shippingRepo) ShipOrder(ctx context.Context, o *biz.ShipOrder) (err error) {
	se := &ShippingEntry{
		OrderId: fmt.Sprintf("%d", o.Id),
	}
	b, err := json.Marshal(se)
	if err != nil {
		return err
	}
	uc.data.kp.Input() <- &sarama.ProducerMessage{
		Topic: "shipping",
		Value: sarama.ByteEncoder(b),
	}
	return
}
