package data

import (
	"github.com/go-kratos/kratos/v2/log"
)

// Data .
type Data struct {
	log *log.Helper
}

// NewData .
func NewData(logger log.Logger) (*Data, error) {
	log := log.NewHelper(log.With(logger, "module", "payment-service/data"))

	d := &Data{
		log: log,
	}
	return d, nil
}
