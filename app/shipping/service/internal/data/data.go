package data

import (
	"github.com/IBM/sarama"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/do/v2"

	"github.com/go-kratos/beer-shop/app/shipping/service/internal/conf"

	// init mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// Data .
type Data struct {
	kp  sarama.AsyncProducer
	log *log.Helper
}

// NewData .
func NewData(i do.Injector) (*Data, error) {
	confData := do.MustInvoke[*conf.Data](i)
	logger := do.MustInvoke[log.Logger](i)
	producer := NewKafkaProducer(confData)
	helper := log.NewHelper(log.With(logger, "module", "shipping-service/data"))
	d := &Data{
		kp:  producer,
		log: helper,
	}
	return d, nil
}

func NewKafkaProducer(conf *conf.Data) sarama.AsyncProducer {
	c := sarama.NewConfig()
	p, err := sarama.NewAsyncProducer(conf.Kafka.Addrs, c)
	if err != nil {
		panic(err)
	}
	return p
}

func (d *Data) Shutdown() error {
	return d.kp.Close()
}
