package data

import (
	"backend/application/order/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type orderRepo struct {
	data *Data
	log  *log.Helper
}

func (o *orderRepo) PlaceOrder(ctx context.Context, req *biz.PlaceOrderReq) (*biz.PlaceOrderResp, error) {

}

func (o *orderRepo) ListOrder(ctx context.Context, req *biz.ListOrderReq) (*biz.ListOrderResp, error) {
	// TODO implement me
	panic("implement me")
}

func (o *orderRepo) MarkOrderPaid(ctx context.Context, req *biz.MarkOrderPaidReq) (*biz.MarkOrderPaidResp, error) {
	// TODO implement me
	panic("implement me")
}

type OrderRepo interface {
	PlaceOrder(ctx context.Context, req *biz.PlaceOrderReq) (*biz.PlaceOrderResp, error)
	ListOrder(ctx context.Context, req *biz.ListOrderReq) (*biz.ListOrderResp, error)
	MarkOrderPaid(ctx context.Context, req *biz.MarkOrderPaidReq) (*biz.MarkOrderPaidResp, error)
}

func NewOrderRepo(data *Data, logger log.Logger) biz.OrderRepo {
	return &orderRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "order/data")),
	}
}
