package service

import (
	v1 "github.com/go-kratos/beer-shop/api/_gen/go/shop/interface/v1"
	"github.com/go-kratos/beer-shop/app/shop/interface/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/do/v2"
)

type ShopInterface struct {
	v1.UnimplementedShopInterfaceServer

	uc *biz.UserUseCase
	ac *biz.AuthUseCase
	cc *biz.CatalogUseCase

	log *log.Helper
}

func NewShopInterface(i do.Injector) (*ShopInterface, error) {
	uc := do.MustInvoke[*biz.UserUseCase](i)
	cc := do.MustInvoke[*biz.CatalogUseCase](i)
	ac := do.MustInvoke[*biz.AuthUseCase](i)
	logger := do.MustInvoke[log.Logger](i)

	return &ShopInterface{
		log: log.NewHelper(log.With(logger, "module", "service/interface")),
		uc:  uc,
		ac:  ac,
		cc:  cc,
	}, nil
}
