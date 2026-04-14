package service

import (
	v1 "github.com/go-kratos/beer-shop/api/_gen/go/catalog/service/v1"
	"github.com/go-kratos/beer-shop/app/catalog/service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type CatalogService struct {
	v1.UnimplementedCatalogServer

	bc  *biz.BeerUseCase
	log *log.Helper
}

func NewCatalogService(bc *biz.BeerUseCase, logger log.Logger) *CatalogService {
	return &CatalogService{

		bc:  bc,
		log: log.NewHelper(log.With(logger, "module", "service/catalog"))}
}
