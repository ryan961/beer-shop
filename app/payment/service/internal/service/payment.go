package service

import (
	"context"

	v1 "github.com/go-kratos/beer-shop/api/_gen/go/payment/service/v1"
)

func (s *PaymentService) PaymentAuth(ctx context.Context, req *v1.PaymentAuthReq) (reply *v1.PaymentAuthReply, err error) {
	return &v1.PaymentAuthReply{}, err
}
