package service

import (
	v1 "github.com/go-kratos/beer-shop/api/_gen/go/shop/admin/v1"
	"github.com/go-kratos/beer-shop/app/shop/admin/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/do/v2"
)

type ShopAdmin struct {
	v1.UnimplementedShopAdminServer

	log *log.Helper
	cc  *biz.CatalogUseCase
	uc  *biz.UserUseCase
}

func NewShopAdmin(i do.Injector) (*ShopAdmin, error) {
	uc := do.MustInvoke[*biz.UserUseCase](i)
	cc := do.MustInvoke[*biz.CatalogUseCase](i)
	logger := do.MustInvoke[log.Logger](i)
	return &ShopAdmin{
		log: log.NewHelper(log.With(logger, "module", "service/interface")),
		uc:  uc,
		cc:  cc,
	}, nil
}
