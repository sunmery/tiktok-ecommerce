package service

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	kerrors "github.com/go-kratos/kratos/v2/errors"

	"backend/constants"

	globalPkg "backend/pkg"

	v1 "backend/api/merchant/order/v1"

	"github.com/go-kratos/kratos/v2/log"

	cartv1 "backend/api/cart/v1"
	orderv1 "backend/api/order/v1"
	userv1 "backend/api/user/v1"

	"google.golang.org/protobuf/types/known/timestamppb"

	"backend/application/merchant/internal/biz"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderService struct {
	v1.UnimplementedOrderServer
	oc *biz.OrderUsecase
}

func NewOrderService(oc *biz.OrderUsecase) *OrderService {
	return &OrderService{oc: oc}
}

// GetMerchantByOrderId 根据订单ID查找商家
func (s *OrderService) GetMerchantByOrderId(ctx context.Context, req *v1.GetMerchantByOrderIdReq) (*v1.GetMerchantByOrderIdReply, error) {
	// 调用业务层获取订单列表
	resp, err := s.oc.GetMerchantByOrderId(ctx, &biz.GetMerchantByOrderIdReq{
		OrderId: req.OrderId,
	})
	if err != nil {
		return nil, err
	}

	return &v1.GetMerchantByOrderIdReply{
		MerchantId: resp.MerchantId.String(),
	}, nil
}

// GetMerchantOrders 获取商家订单列表
func (s *OrderService) GetMerchantOrders(ctx context.Context, req *v1.GetMerchantOrdersReq) (*v1.GetMerchantOrdersReply, error) {
	var userId uuid.UUID
	var err error
	if req.MerchantId == "" {
		// 从网关获取用户ID
		userId, err = globalPkg.GetMetadataUesrID(ctx)
		if err != nil {
			log.Errorf("获取用户ID失败: %v", err)
			return nil, status.Error(codes.Unauthenticated, "无效的商家ID")

		}
	} else {
		// 解析商家ID
		userId, err = uuid.Parse(req.MerchantId)
		if err != nil {
			log.Errorf("解析商家ID失败: %v", err)
			return nil, status.Error(codes.InvalidArgument, "无效的用户ID")
		}
	}

	// 调用业务层获取订单列表
	resp, err := s.oc.GetMerchantOrders(ctx, &biz.GetMerchantOrdersReq{
		UserID:   userId,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	// 检查是否有订单
	if len(resp.Orders) == 0 {
		log.Infof("商家 %s 没有订单记录", userId)
		return &v1.GetMerchantOrdersReply{Orders: nil}, nil
	}

	// 将业务层返回的订单数据转换为proto消息格式
	merchantOrders := make([]*v1.MerchantOrder, 0, len(resp.Orders))
	for _, order := range resp.Orders {
		// 创建订单项
		orderItems := make([]*v1.OrderItem, 0, len(order.Items))
		for _, item := range order.Items {
			// 创建购物车商品
			cartItem := &cartv1.CartItem{
				MerchantId: item.Item.MerchantId.String(),
				ProductId:  item.Item.ProductId.String(),
				Quantity:   item.Item.Quantity,
			}

			// 创建订单项
			orderItem := &v1.OrderItem{
				SubOrderId:     item.SubOrderID,
				Item:           cartItem,
				Cost:           item.Cost,
				Email:          item.Email,
				UserId:         item.UserID.String(),
				Currency:       item.Currency,
				PaymentStatus:  convertToPaymentStatus(item.PaymentStatus),
				ShippingStatus: convertToShippingStatus(item.ShippingStatus),
				CreatedAt:      timestamppb.New(item.CreatedAt),
				UpdatedAt:      timestamppb.New(item.UpdatedAt),
			}

			// 添加地址信息
			if item.Address.StreetAddress != "" {
				orderItem.Address = &userv1.ConsumerAddress{
					StreetAddress: item.Address.StreetAddress,
					City:          item.Address.City,
					State:         item.Address.State,
					Country:       item.Address.Country,
					ZipCode:       item.Address.ZipCode,
				}
			}

			orderItems = append(orderItems, orderItem)
		}

		// 创建商家订单
		merchantOrder := &v1.MerchantOrder{
			Items:     orderItems,
			OrderId:   order.OrderID,
			CreatedAt: timestamppb.New(order.CreatedAt),
		}

		merchantOrders = append(merchantOrders, merchantOrder)
	}

	return &v1.GetMerchantOrdersReply{
		Orders: merchantOrders,
	}, nil
}

// CreateOrderShip 创建货运信息
func (s *OrderService) CreateOrderShip(ctx context.Context, req *v1.CreateOrderShipReq) (*v1.CreateOrderShipReply, error) {
	// 从网关获取用户ID
	userId, err := globalPkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "failed to get user ID")
	}
	// 验证订单ID
	if req.SubOrderId == 0 {
		return nil, status.Error(codes.InvalidArgument, "order ID is required")
	}

	shippingAddress, err := json.Marshal(req.ShippingAddress)
	if err != nil {
		return nil, kerrors.New(500, "shipping_address", err.Error())
	}

	result, err := s.oc.CreateOrderShip(ctx, &biz.CreateOrderShipReq{
		MerchantID:      userId,
		SubOrderId:      req.SubOrderId,
		TrackingNumber:  req.TrackingNumber,
		Carrier:         req.Carrier,
		ShippingStatus:  constants.ShippingShipped,
		ShippingAddress: shippingAddress,
		ShippingFee:     req.ShippingFee,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateOrderShipReply{
		Id:        result.Id,
		CreatedAt: timestamppb.New(result.CreatedAt),
	}, nil
}

func (s *OrderService) UpdateOrderShippingStatus(ctx context.Context, req *v1.UpdateOrderShippingStatusReq) (*v1.UpdateOrderShippingStatusReply, error) {
	// 从网关获取用户ID
	userId, err := globalPkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "failed to get user ID")
	}
	// 验证订单ID
	if req.SubOrderId == 0 {
		return nil, status.Error(codes.InvalidArgument, "order ID is required")
	}
	// 调用业务层获取订单状态
	shippingAddress, err := json.Marshal(req.ShippingAddress)
	if err != nil {
		return nil, kerrors.New(500, "shipping_address", err.Error())
	}

	orderStatus, err := s.oc.UpdateOrderShippingStatus(ctx, &biz.UpdateOrderShippingStatusReq{
		MerchantID:     userId,
		SubOrderId:     req.SubOrderId,
		TrackingNumber: req.TrackingNumber,
		Carrier:        req.Carrier,
		ShippingStatus: convertShippingStatusProtoToBiz(req.ShippingStatus),
		// Delivery:        req.Delivery,
		ShippingAddress: shippingAddress,
		ShippingFee:     req.ShippingFee,
	})
	if err != nil {
		return nil, err
	}
	return &v1.UpdateOrderShippingStatusReply{
		Id:        orderStatus.ID,
		UpdatedAt: timestamppb.New(orderStatus.UpdatedAt),
	}, nil
}

func convertShippingStatusProtoToBiz(status orderv1.ShippingStatus) constants.ShippingStatus {
	switch status {
	case orderv1.ShippingStatus_WAIT_COMMAND:
		return constants.ShippingWaitCommand
	case orderv1.ShippingStatus_PENDING_SHIPMENT:
		return constants.ShippingPending
	case orderv1.ShippingStatus_SHIPPED:
		return constants.ShippingShipped
	case orderv1.ShippingStatus_IN_TRANSIT:
		return constants.ShippingInTransit
	case orderv1.ShippingStatus_DELIVERED:
		return constants.ShippingDelivered
	case orderv1.ShippingStatus_CONFIRMED:
		return constants.ShippingConfirmed
	case orderv1.ShippingStatus_CANCELLED_SHIPMENT:
		return constants.ShippingCancelled
	default:
		return constants.ShippingWaitCommand
	}
}

func convertToShippingStatus(status constants.ShippingStatus) orderv1.ShippingStatus {
	switch status {
	case constants.ShippingWaitCommand:
		return orderv1.ShippingStatus_WAIT_COMMAND
	case constants.ShippingPending:
		return orderv1.ShippingStatus_PENDING_SHIPMENT
	case constants.ShippingShipped:
		return orderv1.ShippingStatus_SHIPPED
	case constants.ShippingInTransit:
		return orderv1.ShippingStatus_IN_TRANSIT
	case constants.ShippingDelivered:
		return orderv1.ShippingStatus_DELIVERED
	case constants.ShippingConfirmed:
		return orderv1.ShippingStatus_CONFIRMED
	case constants.ShippingCancelled:
		return orderv1.ShippingStatus_CANCELLED_SHIPMENT
	default:
		return orderv1.ShippingStatus_WAIT_COMMAND
	}
}

func convertToPaymentStatus(status constants.PaymentStatus) orderv1.PaymentStatus {
	switch status {
	case constants.PaymentPending:
		return orderv1.PaymentStatus_PENDING
	case constants.PaymentPaid:
		return orderv1.PaymentStatus_PAID
	case constants.PaymentFailed:
		return orderv1.PaymentStatus_FAILED
	case constants.PaymentCancelled:
		return orderv1.PaymentStatus_CANCELLED
	default:
		return orderv1.PaymentStatus_PENDING
	}
}
