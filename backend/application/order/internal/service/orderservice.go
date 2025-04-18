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
				Name:       item.Item.Name,
				Picture:    item.Item.Picture,
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

// GetOrder 查询用户订单ID
func (s *OrderServiceService) GetOrder(ctx context.Context, req *v1.GetOrderReq) (*v1.Order, error) {
	userId, err := pkg.GetMetadataUesrID(ctx)
	order, err := s.uc.GetOrder(ctx, &biz.GetOrderReq{
		UserId:  userId,
		OrderId: req.Id,
	})
	if err != nil {
		return nil, err
	}

	return &v1.Order{
		Items:         order.Items,
		OrderId:       order.OrderId,
		UserId:        order.UserId,
		Currency:      order.Currency,
		Address:       order.Address,
		Email:         order.Email,
		CreatedAt:     order.CreatedAt,
		PaymentStatus: order.PaymentStatus,
	}, nil
}

func (s *OrderServiceService) GetConsumerOrders(ctx context.Context, req *v1.GetConsumerOrdersReq) (*v1.Orders, error) {
	// 从网关获取用户ID
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		log.Errorf("获取用户ID失败: %v", err)
		return nil, status.Error(codes.Unauthenticated, "获取用户ID失败")
	}

	// 构建业务层请求
	listReq := &biz.GetConsumerOrdersReq{
		UserId:   userId,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	// 调用业务层获取订单列表
	resp, err := s.uc.GetConsumerOrders(ctx, listReq)
	if err != nil {
		log.Errorf("获取用户订单失败: %v", err)
		return nil, status.Errorf(codes.Internal, "获取用户订单失败: %v", err)
	}

	// 直接返回业务层的响应，因为它已经是v1.Orders的格式
	if len(resp.Orders) == 0 {
		log.Infof("用户 %s 没有订单记录", userId)
		return &v1.Orders{Orders: []*v1.Order{}}, nil
	}

	return &v1.Orders{Orders: resp.Orders}, nil
}

func (s *OrderServiceService) GetAllOrders(ctx context.Context, req *v1.GetAllOrdersReq) (*v1.Orders, error) {
	// 从网关获取用户ID
	userId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		log.Errorf("获取用户ID失败: %v", err)
		return nil, status.Error(codes.Unauthenticated, "获取用户ID失败")
	}

	// 使用请求中的merchant_id覆盖，如果有指定的话

	// 构建业务层请求
	listReq := &biz.GetAllOrdersReq{
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	// 调用业务层获取订单列表
	resp, err := s.uc.GetAllOrders(ctx, listReq)
	if err != nil {
		log.Errorf("获取用户订单失败: %v", err)
		return nil, status.Errorf(codes.Internal, "获取用户订单失败: %v", err)
	}

	// 检查是否有订单
	if len(resp.Orders) == 0 {
		log.Infof("用户 %s 没有订单记录", userId)
		return &v1.Orders{Orders: []*v1.Order{}}, nil
	}

	// 按照用户订单分组
	merchantOrders := make(map[int64][]*biz.SubOrder)
	for _, subOrder := range resp.Orders {
		merchantOrders[subOrder.ID] = append(merchantOrders[subOrder.ID], subOrder)
	}

	// 转换订单列表为API响应格式
	var orders []*v1.Order
	for _, subOrders := range merchantOrders {
		if len(subOrders) == 0 {
			continue
		}

		// 使用第一个子订单信息
		firstSubOrder := subOrders[0]

		// 订单项集合 - 汇总所有子订单的订单项
		var orderItems []*v1.OrderItem
		for _, subOrder := range subOrders {
			for _, item := range subOrder.Items {
				// 确保CartItem中的数据是有效的
				if item.Item == nil {
					log.Warnf("跳过缺少商品信息的订单项, 订单ID: %d", subOrder.ID)
					continue
				}

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
		createdAt := timestamppb.New(firstSubOrder.CreatedAt)

		// 解析支付状态
		paymentStatus := mapStringStatusToProto(firstSubOrder.Status)

		// 创建地址信息 (在真实场景中需要从订单数据中获取)
		address := &userv1.Address{
			StreetAddress: "未提供地址信息", // 这里应该从订单数据中获取实际地址
			City:          "",
			State:         "",
			Country:       "",
			ZipCode:       "",
		}

		// 添加订单到响应列表
		orders = append(orders, &v1.Order{
			Items:         orderItems,
			OrderId:       firstSubOrder.ID, // 注意: 确保ID类型转换正确
			UserId:        firstSubOrder.MerchantID.String(),
			Currency:      firstSubOrder.Currency,
			Address:       address,
			Email:         "未提供邮箱", // 这里应该从订单数据中获取实际邮箱
			CreatedAt:     createdAt,
			PaymentStatus: paymentStatus,
		})
	}

	log.Debugf("返回 %d 个用户订单", len(orders))
	return &v1.Orders{Orders: orders}, nil
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

// 将字符串状态转换为Proto枚举
func mapStringStatusToProto(status string) v1.PaymentStatus {
	switch status {
	case string(biz.PaymentPending):
		return v1.PaymentStatus_NOT_PAID
	case string(biz.PaymentProcessing):
		return v1.PaymentStatus_PROCESSING
	case string(biz.PaymentPaid):
		return v1.PaymentStatus_PAID
	case string(biz.PaymentFailed):
		return v1.PaymentStatus_FAILED
	case string(biz.PaymentCancelled):
		return v1.PaymentStatus_CANCELLED
	default:
		return v1.PaymentStatus_NOT_PAID
	}
}
