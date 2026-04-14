package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/go-kratos/beer-shop/app/cart/service/internal/conf"
)

// Data .
type Data struct {
	db  *mongo.Database
	log *log.Helper
}

func NewMongo(conf *conf.Data) *mongo.Database {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.Mongodb.Uri))
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}
	return client.Database(conf.Mongodb.Database)
}

// NewData .
func NewData(database *mongo.Database, logger log.Logger) (*Data, error) {
	log := log.NewHelper(log.With(logger, "module", "cart-service/data"))

	d := &Data{
		db:  database,
		log: log,
	}
	return d, nil
}

func (d *Data) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return d.db.Client().Disconnect(ctx)
}
