package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type (
	CheckoutRequest struct {
		UserId        uuid.UUID // 用户 ID（可选），如果用户未注册，则可以为空
		Firstname     string    // 用户的名字（必填），允许非注册用户直接填写信息下单
		Lastname      string    // 用户的姓氏（必填）
		Email         string    // 用户的邮箱地址（必填），用于接收订单确认邮件等
		CreditCardId  uint32    // 用户的信用卡信息（必填），用于支付
		AddressId     uint32
		PaymentMethod string
	}
	CheckoutReply struct {
		OrderId    string // 唯一标识订单，用于后续查询、退换货等操作
		PaymentId  string // 支付事务唯一标识，用于与支付网关对账或处理支付相关问题
		PaymentURL string // 支付URL
	}
)

type CheckoutRepo interface {
	Checkout(context.Context, *CheckoutRequest) (*CheckoutReply, error)
}

type CheckoutUsecase struct {
	repo CheckoutRepo
	log  *log.Helper
}

func NewCheckoutUsecase(repo CheckoutRepo, logger log.Logger) *CheckoutUsecase {
	return &CheckoutUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (cc *CheckoutUsecase) Checkout(ctx context.Context, req *CheckoutRequest) (*CheckoutReply, error) {
	cc.log.WithContext(ctx).Debugf("Signin request: %+v", req)
	return cc.repo.Checkout(ctx, req)
}
