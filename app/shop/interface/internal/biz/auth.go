package biz

import (
	"context"
	"errors"

	v1 "github.com/go-kratos/beer-shop/api/_gen/go/shop/interface/v1"
	"github.com/go-kratos/beer-shop/app/shop/interface/internal/conf"

	"github.com/golang-jwt/jwt"
	"github.com/samber/do/v2"
)

var (
	ErrLoginFailed = errors.New("login failed")
)

type AuthUseCase struct {
	key      string
	userRepo UserRepo
}

func NewAuthUseCase(i do.Injector) (*AuthUseCase, error) {
	authConfig := do.MustInvoke[*conf.Auth](i)
	userRepo := do.MustInvoke[UserRepo](i)

	return &AuthUseCase{
		key:      authConfig.ApiKey,
		userRepo: userRepo,
	}, nil
}

func (receiver *AuthUseCase) Login(ctx context.Context, req *v1.LoginReq) (*v1.LoginReply, error) {

	// get user
	user, err := receiver.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, v1.ErrorLoginFailed("user not found: %s", err.Error())
	}
	// check permission(password blacklist etc...)
	err = receiver.userRepo.VerifyPassword(ctx, user, req.Password)
	if err != nil {
		return nil, v1.ErrorLoginFailed("password not match")
	}
	// generate token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
	})
	signedString, err := claims.SignedString([]byte(receiver.key))
	if err != nil {
		return nil, v1.ErrorLoginFailed("generate token failed: %s", err.Error())
	}
	return &v1.LoginReply{
		Token: signedString,
	}, nil
}

func (receiver *AuthUseCase) Register(ctx context.Context, req *v1.RegisterReq) (*v1.RegisterReply, error) {

	// check username
	_, err := receiver.userRepo.FindByUsername(ctx, req.Username)
	if !errors.Is(err, ErrUserNotFound) {
		return nil, v1.ErrorRegisterFailed("username already exists")
	}
	// create user
	user, err := NewUser(req.Username, req.Password)
	if err != nil {
		return nil, v1.ErrorRegisterFailed("create user failed: %s", err.Error())
	}
	// save user
	err = receiver.userRepo.Save(ctx, &user)
	if err != nil {
		return nil, v1.ErrorRegisterFailed("save user failed: %s", err.Error())
	}
	return &v1.RegisterReply{
		Id: user.Id,
	}, nil
}
