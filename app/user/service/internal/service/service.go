package service

import (
	v1 "github.com/go-kratos/beer-shop/api/_gen/go/user/service/v1"
	"github.com/go-kratos/beer-shop/app/user/service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/do/v2"
)

type UserService struct {
	v1.UnimplementedUserServer

	uc  *biz.UserUseCase
	ac  *biz.AddressUseCase
	cc  *biz.CardUseCase
	log *log.Helper
}

func NewUserService(i do.Injector) (*UserService, error) {
	uc := do.MustInvoke[*biz.UserUseCase](i)
	cc := do.MustInvoke[*biz.CardUseCase](i)
	ac := do.MustInvoke[*biz.AddressUseCase](i)
	logger := do.MustInvoke[log.Logger](i)
	return &UserService{
		uc:  uc,
		ac:  ac,
		cc:  cc,
		log: log.NewHelper(log.With(logger, "module", "service/server-service")),
	}, nil
}
