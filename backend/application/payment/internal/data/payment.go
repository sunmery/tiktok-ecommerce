package data

import (
	"backend/application/payment/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type PaymentRepo struct {
	data *Data
	log  *log.Helper
}

func NewPaymentRepo(data *Data, logger log.Logger) biz.PaymentRepo {
	return &PaymentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *PaymentRepo) Create(ctx context.Context, req *biz.CreateRequest) (*biz.CreateReply, error) {
	r.log.Infof("Create request: %+v", req)
	return &biz.CreateReply{}, nil
}
