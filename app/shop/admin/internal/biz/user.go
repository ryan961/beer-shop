package biz

import (
	usV1 "github.com/go-kratos/beer-shop/api/_gen/go/user/service/v1"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/do/v2"
)

type User struct {
	Id int64
}

type UserRepo interface {
}

type UserUseCase struct {
	repo UserRepo
	us   usV1.UserClient
	log  *log.Helper
}

func NewUserUseCase(i do.Injector) (*UserUseCase, error) {
	repo := do.MustInvoke[UserRepo](i)
	logger := do.MustInvoke[log.Logger](i)
	us := do.MustInvoke[usV1.UserClient](i)
	helper := log.NewHelper(log.With(logger, "module", "usecase/interface"))
	return &UserUseCase{
		repo: repo,
		us:   us,
		log:  helper,
	}, nil
}
