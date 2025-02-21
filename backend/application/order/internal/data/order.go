package data

import (
	"backend/application/order/internal/biz"
	"backend/application/order/internal/data/models"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
)

type orderRepo struct {
	data *Data
	log  *log.Helper
}

func (o *orderRepo) PlaceOrder(ctx context.Context, req *biz.PlaceOrderReq) (*biz.PlaceOrderResp, error) {
	order, err := o.data.DB(ctx).CreateOrder(ctx, models.CreateOrderParams{
		UserID:        req.UserId,
		Currency:      req.Currency,
		StreetAddress: req.Address.StreetAddress,
		City:          req.Address.City,
		State:         req.Address.State,
		Country:       req.Address.Country,
		ZipCode:       int32(req.Address.ZipCode),
		Email:         req.Email,
	})
	fmt.Printf("order: %v", order)
	if err != nil {
		return nil, err
	}
	return &biz.PlaceOrderResp{
		Order: &biz.OrderResult{
			OrderId: order.ID,
		},
	}, nil
}

func (o *orderRepo) ListOrder(ctx context.Context, req *biz.ListOrderReq) (*biz.ListOrderResp, error) {
	// TODO implement me
	panic("implement me")
}

func (o *orderRepo) MarkOrderPaid(ctx context.Context, req *biz.MarkOrderPaidReq) (*biz.MarkOrderPaidResp, error) {
	// TODO implement me
	panic("implement me")
}

func NewOrderRepo(data *Data, logger log.Logger) biz.OrderRepo {
	return &orderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
