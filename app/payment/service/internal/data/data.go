package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/do/v2"
)

// Data .
type Data struct {
	log *log.Helper
}

// NewData .
func NewData(i do.Injector) (*Data, error) {
	logger := do.MustInvoke[log.Logger](i)
	helper := log.NewHelper(log.With(logger, "module", "payment-service/data"))

	d := &Data{
		log: helper,
	}
	return d, nil
}
