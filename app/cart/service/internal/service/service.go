package service

import (
	v1 "github.com/go-kratos/beer-shop/api/_gen/go/cart/service/v1"
	"github.com/go-kratos/beer-shop/app/cart/service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type CartService struct {
	v1.UnimplementedCartServer

	cc  *biz.CartUseCase
	log *log.Helper
}

func NewCartService(cc *biz.CartUseCase, logger log.Logger) *CartService {
	return &CartService{
		cc:  cc,
		log: log.NewHelper(log.With(logger, "module", "service/cart"))}
}
