package data

import (
	"context"
	"github.com/go-kratos/beer-shop/app/shop/admin/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/do/v2"

	userv1 "github.com/go-kratos/beer-shop/api/_gen/go/user/service/v1"
)

var _ biz.UserRepo = (*userRepo)(nil)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(i do.Injector) (biz.UserRepo, error) {
	data := do.MustInvoke[*Data](i)
	logger := do.MustInvoke[log.Logger](i)
	return &userRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/user")),
	}, nil
}


func (r *userRepo) GetUser(ctx context.Context, id int64) (*biz.User, error) {
	reply, err := r.data.uc.GetUser(ctx, &userv1.GetUserReq{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return &biz.User{
		Id:          reply.Id,
	}, err
}

func (r *userRepo) ListUser(ctx context.Context, pageNum, pageSize int64) ([]*biz.User, error) {
	reply, err := r.data.uc.ListUser(ctx, &userv1.ListUserReq{
	})
	if err != nil {
		return nil, err
	}
	rv := make([]*biz.User, 0)
	for _, x := range reply.Results {
		rv = append(rv, &biz.User{
			Id:          x.Id,
		})
	}
	return rv, err
}
