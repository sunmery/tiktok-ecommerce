package service

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"

	cartv1 "backend/api/cart/v1"
	userv1 "backend/api/user/v1"

	"google.golang.org/protobuf/types/known/timestamppb"

	"backend/application/order/internal/biz"
	"backend/pkg"

	v1 "backend/api/order/v1"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderServiceService struct {
	v1.UnimplementedOrderServiceServer

	uc *biz.OrderUsecase
}

func NewOrderServiceService(uc *biz.OrderUsecase) *OrderServiceService {
	return &OrderServiceService{uc: uc}
}

func (s *OrderServiceService) PlaceOrder(ctx context.Context, req *v1.PlaceOrderReq) (*v1.PlaceOrderResp, error) {
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
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
				Quantity:   item.Item.Quantity,
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
		// Url: order.URL,
	}, nil
}

func (s *OrderServiceService) QueryOrders(ctx context.Context, req *v1.ListOrderReq) (*v1.ListOrderResp, error) {
	// 从网关获取用户ID
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "failed to get user ID")
	}

	// 构建业务层请求
	listReq := &biz.ListOrderReq{
		UserID:   userId,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	// 调用业务层获取订单列表
	resp, err := s.uc.ListOrder(ctx, listReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list orders: %v", err)
	}

	// 转换订单列表为API响应格式
	var orders []*v1.Order
	for _, order := range resp.Orders {
		// 转换地址信息
		address := &userv1.Address{
			StreetAddress: order.Address.StreetAddress,
			City:          order.Address.City,
			State:         order.Address.State,
			Country:       order.Address.Country,
			ZipCode:       order.Address.ZipCode,
		}

		// 转换订单项
		var orderItems []*v1.OrderItem
		for _, subOrder := range order.SubOrders {
			for _, item := range subOrder.Items {
				orderItems = append(orderItems, &v1.OrderItem{
					Item: &cartv1.CartItem{
						MerchantId: item.Item.MerchantId.String(),
						ProductId:  item.Item.ProductId.String(),
						Quantity:   item.Item.Quantity,
					},
					Cost: item.Cost,
				})
			}
		}

		// 转换时间戳
		createdAt := timestamppb.New(order.CreatedAt)

		// 添加订单到响应列表
		orders = append(orders, &v1.Order{
			Items:         orderItems,
			OrderId:       order.OrderID,
			UserId:        order.UserID.String(),
			Currency:      order.Currency,
			Address:       address,
			Email:         order.Email,
			CreatedAt:     createdAt,
			PaymentStatus: mapBizStatusToProto(order.PaymentStatus),
		})
	}

	return &v1.ListOrderResp{Orders: orders}, nil
}

func (s *OrderServiceService) MarkOrderPaid(ctx context.Context, req *v1.MarkOrderPaidReq) (*v1.MarkOrderPaidResp, error) {
	// 从网关获取用户ID
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "failed to get user ID")
	}

	// 验证订单ID
	if req.OrderId == 0 {
		return nil, status.Error(codes.InvalidArgument, "order ID is required")
	}

	// 调用业务层标记订单为已支付
	orderPaid, err := s.uc.MarkOrderPaid(ctx, &biz.MarkOrderPaidReq{
		UserId:  userId,
		OrderId: req.OrderId,
	})
	if err != nil {
		// 根据错误类型返回不同的状态码
		if err.Error() == "order does not belong to user" {
			return nil, status.Error(codes.PermissionDenied, "order does not belong to user")
		} else if err.Error() == "invalid order ID format" {
			return nil, status.Error(codes.InvalidArgument, "invalid order ID format")
		} else if err.Error() == "failed to get order" {
			return nil, status.Error(codes.NotFound, "order not found")
		}
		// 其他错误作为内部错误处理
		return nil, status.Errorf(codes.Internal, "failed to mark order as paid: %v", err)
	}

	log.Debugf("orderPaid: %v", orderPaid)

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
