package service

import (
	v1 "github.com/go-kratos/beer-shop/api/_gen/go/shop/admin/v1"
	"github.com/go-kratos/beer-shop/app/shop/admin/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type ShopAdmin struct {
	v1.UnimplementedShopAdminServer

	log *log.Helper
	cc  *biz.CatalogUseCase
	uc  *biz.UserUseCase
}

func NewShopAdmin(uc *biz.UserUseCase, cc *biz.CatalogUseCase, logger log.Logger) *ShopAdmin {
	return &ShopAdmin{
		log: log.NewHelper(log.With(logger, "module", "service/interface")),
		uc:  uc,
		cc:  cc,
	}
}
