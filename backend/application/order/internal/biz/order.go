package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/google/uuid"
)

type (
	MarkOrderPaidReq struct {
		UserId  uuid.UUID
		OrderId int64
	}
	MarkOrderPaidResp struct{}
)

type OrderRepo interface {
	MarkOrderPaid(ctx context.Context, req *MarkOrderPaidReq) (*MarkOrderPaidResp, error)
}
type OrderUsecase struct {
	repo OrderRepo
	log  *log.Helper
}

func NewUserUsecase(repo OrderRepo, logger log.Logger) *OrderUsecase {
	return &OrderUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (oc *OrderUsecase) MarkOrderPaid(ctx context.Context, req *MarkOrderPaidReq) (*MarkOrderPaidResp, error) {
	oc.log.WithContext(ctx).Debugf("biz/order MarkOrderPaid req:%+v", req)
	return oc.repo.MarkOrderPaid(ctx, req)
}
