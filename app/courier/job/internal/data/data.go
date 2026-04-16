package data

import (
	"context"
	"errors"

	orderv1 "github.com/go-kratos/beer-shop/api/_gen/go/order/service/v1"
	"github.com/go-kratos/beer-shop/app/courier/job/internal/conf"

	"github.com/IBM/sarama"
	consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	consulAPI "github.com/hashicorp/consul/api"
	"github.com/samber/do/v2"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	// init mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// Data .
type Data struct {
	kc  sarama.Consumer
	oc  orderv1.OrderClient
	log *log.Helper
}

// NewData .
func NewData(i do.Injector) (*Data, error) {
	confData := do.MustInvoke[*conf.Data](i)
	logger := do.MustInvoke[log.Logger](i)
	consumer := NewKafkaConsumer(confData)
	oc := do.MustInvoke[orderv1.OrderClient](i)
	helper := log.NewHelper(log.With(logger, "module", "courier-job/data"))
	d := &Data{
		kc:  consumer,
		oc:  oc,
		log: helper,
	}
	return d, nil
}

func NewKafkaConsumer(conf *conf.Data) sarama.Consumer {
	c := sarama.NewConfig()
	p, err := sarama.NewConsumer(conf.Kafka.Addrs, c)
	if err != nil {
		panic(err)
	}
	return p
}

func NewDiscovery(i do.Injector) (registry.Discovery, error) {
	conf := do.MustInvoke[*conf.Registry](i)
	if conf == nil || conf.Consul == nil {
		return nil, errors.New("registry consul config is required")
	}
	c := consulAPI.DefaultConfig()
	c.Address = conf.Consul.Address
	c.Scheme = conf.Consul.Scheme
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		return nil, err
	}
	r := consul.New(cli, consul.WithHealthCheck(false))
	return r, nil
}

func NewOrderServiceClient(i do.Injector) (orderv1.OrderClient, error) {
	r := do.MustInvoke[registry.Discovery](i)
	tp := do.MustInvoke[*tracesdk.TracerProvider](i)
	if r == nil {
		return nil, errors.New("service discovery is required")
	}
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///beer.order.service"),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			tracing.Client(tracing.WithTracerProvider(tp)),
			recovery.Recovery(),
		),
	)
	if err != nil {
		return nil, err
	}
	return orderv1.NewOrderClient(conn), nil
}

func (d *Data) Shutdown() error {
	return d.kc.Close()
}
