package data

import (
	"context"
	"github.com/go-kratos/beer-shop/app/catalog/service/internal/data/ent/migrate"

	"github.com/go-kratos/beer-shop/app/catalog/service/internal/data/ent"
	"github.com/go-kratos/kratos/v2/log"

	"github.com/go-kratos/beer-shop/app/catalog/service/internal/conf"

	// init mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// Data .
type Data struct {
	db  *ent.Client
	log *log.Helper
}

func NewEntClient(conf *conf.Data, logger log.Logger) *ent.Client {
	log := log.NewHelper(log.With(logger, "module", "catalog-service/data/ent"))

	client, err := ent.Open(
		conf.Database.Driver,
		conf.Database.Source,
	)
	if err != nil {
		log.Fatalf("failed opening connection to db: %v", err)
	}
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background(), migrate.WithForeignKeys(false)); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return client
}

// NewData .
func NewData(entClient *ent.Client, logger log.Logger) (*Data, error) {
	log := log.NewHelper(log.With(logger, "module", "catalog-service/data"))

	d := &Data{
		db:  entClient,
		log: log,
	}
	return d, nil
}

func (d *Data) Shutdown() error {
	return d.db.Close()
}
