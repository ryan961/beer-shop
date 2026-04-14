package data

import (
	"github.com/Shopify/sarama"
	"github.com/go-kratos/kratos/v2/log"

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
func NewData(producer sarama.AsyncProducer, conf *conf.Data, logger log.Logger) (*Data, error) {
	log := log.NewHelper(log.With(logger, "module", "shipping-service/data"))
	d := &Data{
		kp:  producer,
		log: log,
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
