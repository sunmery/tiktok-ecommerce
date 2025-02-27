package biz

import (
	orderv1 "backend/api/order/v1"
	"context"
	"github.com/google/uuid"
	"time"
)

const (
	SystemMessageTemplate = `你是一个智能助手，需要处理以下类型的请求（当前时间：%s）：
1. 当用户要求购买商品时（如"帮我下单商品A"），回复格式：[ORDER]商品名称
2. 当用户要求查看订单时（如"我的订单"），回复格式：[QUERY]
3. 其他情况正常聊天`

	UserMessageTemplate = "用户问题：%s"
)

type AssistantUsecase struct {
	orderRepo OrderRepo
}

type (
	OrderReq struct {
		UserID  uuid.UUID
		MerchantId uuid.UUID
		Product string
	}
	OrderResponse struct {
		ID        string
		Product   string
		Status    string
		CreatedAt time.Time
	}
)

type OrderRepo interface {
	CreateOrder(ctx context.Context, req *OrderReq) (*orderv1.PlaceOrderResp, error)
	QueryOrders(ctx context.Context, req *orderv1.ListOrderReq) (*orderv1.ListOrderResp, error)
	QueryOrders(ctx context.Context, req *orderv1.ListOrderReq) (*orderv1.ListOrderResp, error)
}

func NewAssistantUsecase(repo OrderRepo) *AssistantUsecase {
	return &AssistantUsecase{orderRepo: repo}
}

func (uc *AssistantUsecase) Process(ctx context.Context, question string) (interface{}, error) {
	// TODO implement me
	// 先判断要查询的商品是否存在库存, 如果存在, 那么就创建一个订单
	//
	panic("implement me")
}
