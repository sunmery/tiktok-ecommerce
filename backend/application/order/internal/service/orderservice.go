package service

import (
	"backend/application/order/internal/biz"
	"backend/pkg"
	"context"
	"fmt"

	v1 "backend/api/order/v1"
	"github.com/google/uuid"
)

type OrderServiceService struct {
	v1.UnimplementedOrderServiceServer

	uc *biz.OrderUsecase
}

func NewOrderServiceService(uc *biz.OrderUsecase) *OrderServiceService {
	return &OrderServiceService{uc: uc}
}

func (s *OrderServiceService) PlaceOrder(ctx context.Context, req *v1.PlaceOrderReq) (*v1.PlaceOrderResp, error) {
	// 从网关获取用户ID

	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}
	var orderItems []*biz.OrderItem
	for _, item := range req.OrderItems {
		merchantId, err := uuid.Parse(item.Item.MerchantId)
		if err != nil {
			return nil, fmt.Errorf("invalid merchant id: %s", item.Item.MerchantId)
		}

		ProductId := uuid.MustParse(item.Item.ProductId)
		orderItems = append(orderItems, &biz.OrderItem{
			Item: &biz.CartItem{
				MerchantId: merchantId,
				ProductId:  ProductId,
				Quantity:   uint32(item.Item.Quantity),
			},
			Cost: item.Cost,
		})
	}

	order, err := s.uc.PlaceOrder(ctx, &biz.PlaceOrderReq{
		UserId:   userId,
		Currency: req.Currency,
		Address: &biz.Address{
			StreetAddress: req.Address.StreetAddress,
			City:          req.Address.City,
			State:         req.Address.State,
			Country:       req.Address.Country,
			ZipCode:       req.Address.ZipCode,
		},
		Email:      req.Email,
		OrderItems: orderItems,
	})
	if err != nil {
		return nil, err
	}

	return &v1.PlaceOrderResp{
		Order: &v1.OrderResult{
			OrderId: order.Order.OrderId,
		},
		Url: order.URL,
	}, nil
}

func (s *OrderServiceService) ListOrders(ctx context.Context, req *v1.ListOrderReq) (*v1.ListOrderResp, error) {
	// 从网关获取用户ID
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}
	listReq := &biz.ListOrderReq{
		UserID:   userId,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	resp, err := s.uc.ListOrder(ctx, listReq)
	if err != nil {
		return nil, err
	}
	fmt.Printf("resp: %+v\n", resp)

	var orders []*v1.Order
	// for _, order := range resp.Orders {
	// 	orders = append(orders, &v1.Order{
	// 		OrderItems:    &v1.OrderItem{
	// 			Item: &v1.CartItem{
	// 				MerchantId: order.OrderItems[0].Item.MerchantId.String(),
	// 				ProductId:  order.OrderItems[0].Item.ProductId.String(),
	// 				Quantity:   uint32(order.OrderItems[0].Item.Quantity),
	// 			},
	// 			Cost: order.Cost,
	// 		},
	// 		OrderId:       order.OrderID,
	// 		UserId:        order.UserID,
	// 		Currency:      order.Currency,
	// 		Address:       nil,
	// 		Email:         order.Email,
	// 		CreatedAt:     order.CreatedAt,
	// 		PaymentStatus: mapBizStatusToProto(order.PaymentStatus),
	// 	})
	// }

	return &v1.ListOrderResp{Orders: orders}, nil
}

func (s *OrderServiceService) MarkOrderPaid(ctx context.Context, req *v1.MarkOrderPaidReq) (*v1.MarkOrderPaidResp, error) {
	// 从网关获取用户ID
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, err
	}
	orderPaid, err := s.uc.MarkOrderPaid(ctx, &biz.MarkOrderPaidReq{
		UserId:  userId,
		OrderId: req.OrderId,
	})
	if err != nil {
		return nil, err
	}
	fmt.Printf("orderPaid:%+v\n", orderPaid)

	return &v1.MarkOrderPaidResp{}, nil
}

// 转换业务层枚举到 Proto int
func mapBizStatusToProto(status biz.PaymentStatus) v1.PaymentStatus {
	switch status {
	case biz.PaymentPending:
		return v1.PaymentStatus_NOT_PAID
	case biz.PaymentProcessing:
		return v1.PaymentStatus_PROCESSING
	case biz.PaymentPaid:
		return v1.PaymentStatus_PAID
	case biz.PaymentFailed:
		return v1.PaymentStatus_FAILED
	case biz.PaymentCancelled:
		return v1.PaymentStatus_CANCELLED
	default:
		return v1.PaymentStatus_NOT_PAID
	}
}
