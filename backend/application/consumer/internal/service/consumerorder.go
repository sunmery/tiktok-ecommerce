package service

import (
	"context"
	"fmt"
	"time"

	"backend/application/consumer/internal/pkg"

	"github.com/go-kratos/kratos/v2/log"

	globalPkg "backend/pkg"

	"backend/application/consumer/internal/biz"

	cartv1 "backend/api/cart/v1"

	userv1 "backend/api/user/v1"
	globalpkg "backend/pkg"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "backend/api/consumer/order/v1"
)

type ConsumerOrderService struct {
	v1.UnimplementedConsumerOrderServer
	oc *biz.ConsumerOrderUsecase
}

func NewConsumerOrderService(oc *biz.ConsumerOrderUsecase) *ConsumerOrderService {
	return &ConsumerOrderService{oc: oc}
}

func (s *ConsumerOrderService) PlaceOrder(ctx context.Context, req *v1.PlaceOrderRequest) (*v1.PlaceOrderReply, error) {
	userId, err := globalpkg.GetMetadataUesrID(ctx)
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

	order, err := s.oc.PlaceOrder(ctx, &biz.PlaceOrderRequest{
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

	return &v1.PlaceOrderReply{
		Order: &v1.PlaceOrderReply_OrderResult{
			OrderId:         order.Order.OrderId,
			FreezeId:        order.Order.FreezeId,
			ConsumerVersion: order.Order.ConsumerVersion,
			MerchantVersion: order.Order.MerchantVersion,
		},
		// Url: order.URL,
	}, nil
}

// GetConsumerOrdersWithSuborders 根据用户主订单查询子订单
func (s *ConsumerOrderService) GetConsumerOrdersWithSuborders(ctx context.Context, req *v1.GetConsumerOrdersWithSubordersRequest) (*v1.GetConsumerOrdersWithSubordersReply, error) {
	var userId uuid.UUID
	var err error
	if req.UserId == "" {
		userId, err = globalpkg.GetMetadataUesrID(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		userId, err = uuid.Parse(req.UserId)
		if err != nil {
			return nil, err
		}
	}

	reply, getUserOrdersWithSubordersErr := s.oc.GetConsumerOrdersWithSuborders(ctx, &biz.GetConsumerOrdersWithSubordersRequest{
		UserId:  userId,
		OrderId: req.OrderId,
	})
	if getUserOrdersWithSubordersErr != nil {
		return nil, getUserOrdersWithSubordersErr
	}

	if reply == nil {
		return nil, nil
	}

	orders := make([]*v1.GetConsumerOrdersWithSubordersReply_Suborders, 0, len(reply.Orders))
	for _, s := range reply.Orders {
		var allItems []*v1.OrderItem
		if s.Items != nil {
			for _, item := range allItems {
				allItems = append(allItems, &v1.OrderItem{
					Item: &cartv1.CartItem{
						MerchantId: item.Item.MerchantId,
						ProductId:  item.Item.ProductId,
						Quantity:   item.Item.Quantity,
						Name:       item.Item.Name,
						Picture:    item.Item.Picture,
					},
					Cost: item.Cost,
				})
			}
		}
		orders = append(orders, &v1.GetConsumerOrdersWithSubordersReply_Suborders{
			Id:             s.OrderId,
			SubOrderId:     s.SubOrderId,
			StreetAddress:  s.StreetAddress,
			City:           s.City,
			State:          s.State,
			Country:        s.Country,
			ZipCode:        s.ZipCode,
			Email:          s.Email,
			MerchantId:     s.MerchantId,
			PaymentStatus:  string(s.PaymentStatus),
			ShippingStatus: string(s.ShippingStatus),
			TotalAmount:    s.TotalAmount,
			Currency:       s.Currency,
			Items:          allItems,
			CreatedAt:      timestamppb.New(s.CreatedAt),
			UpdatedAt:      timestamppb.New(s.UpdatedAt),
		})
	}
	return &v1.GetConsumerOrdersWithSubordersReply{
		Orders: orders,
	}, nil
}

func (s *ConsumerOrderService) GetConsumerOrders(ctx context.Context, req *v1.GetConsumerOrdersRequest) (*v1.ConsumerOrders, error) {
	// 从网关获取用户ID
	var userId uuid.UUID
	var err error
	if req.UserId == "" {
		userId, err = globalpkg.GetMetadataUesrID(ctx)
		if err != nil {
			log.Errorf("获取用户ID失败: %v", err)
			return nil, status.Error(codes.Unauthenticated, "获取用户ID失败")
		}
	} else {
		userId, err = uuid.Parse(req.UserId)
		if err != nil {
			log.Errorf("解析用户ID失败: %v", err)
			return nil, status.Error(codes.InvalidArgument, "解析用户ID失败")
		}
	}

	// 调用业务层获取订单列表
	resp, err := s.oc.GetConsumerOrders(ctx, &biz.GetConsumerOrdersRequest{
		UserId:   userId,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		log.Errorf("获取用户订单失败: %v", err)
		return nil, status.Errorf(codes.Internal, "获取用户订单失败: %v", err)
	}

	if resp == nil {
		log.Infof("用户 %s 没有订单记录", userId)
		return nil, nil
	}
	orders := make([]*v1.Order, 0, len(resp.SubOrders))
	for _, o := range resp.SubOrders {
		items := make([]*v1.OrderItem, 0, len(o.Items))
		for _, item := range o.Items {
			items = append(items, &v1.OrderItem{
				Item: &cartv1.CartItem{
					MerchantId: item.Item.MerchantId.String(),
					ProductId:  item.Item.ProductId.String(),
					Quantity:   item.Item.Quantity,
					Name:       item.Item.Name,
					Picture:    item.Item.Picture,
				},
				Cost: item.Cost,
			})
		}
		orders = append(orders, &v1.Order{
			Items:      items,
			OrderId:    &o.OrderId,
			SubOrderId: &o.SubOrderID,
			UserId:     req.UserId,
			Currency:   o.Currency,
			Address: &userv1.ConsumerAddress{
				UserId:        req.UserId, // 数据库无需存储用户ID, 直接从请求中返回
				City:          o.Address.City,
				State:         o.Address.State,
				Country:       o.Address.Country,
				ZipCode:       o.Address.ZipCode,
				StreetAddress: o.Address.StreetAddress,
			},
			Email:          o.Email,
			CreatedAt:      timestamppb.New(o.CreatedAt),
			PaymentStatus:  pkg.MapPaymentStatusToProto(o.PaymentStatus),
			ShippingStatus: pkg.MapShippingStatusToProto(o.ShippingStatus),
		})
	}

	return &v1.ConsumerOrders{
		Orders: orders,
	}, nil
}

// ConfirmReceived 确认收货(用户角色)
func (s *ConsumerOrderService) ConfirmReceived(ctx context.Context, req *v1.ConfirmReceivedRequest) (*v1.ConfirmReceivedReply, error) {
	// 从网关获取用户ID
	userId, err := globalpkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "failed to get user ID")
	}
	// 验证订单ID
	if req.OrderId == 0 {
		return nil, status.Error(codes.InvalidArgument, "order ID is required")
	}
	// 调用业务层标记订单为已支付
	orderPaid, err := s.oc.ConfirmReceived(ctx, &biz.ConfirmReceivedRequest{
		UserId:  userId,
		OrderId: req.OrderId,
	})
	if err != nil {
		// 根据错误类型返回不同的状态码
		if err.Error() == "order does not belong to user" {
			return nil, status.Error(codes.PermissionDenied, "order does not belong to user")
		}
	}
	log.Debugf("orderPaid: %v", orderPaid)
	return &v1.ConfirmReceivedReply{}, nil
}

func (s *ConsumerOrderService) GetShipOrderStatus(ctx context.Context, req *v1.GetShipOrderStatusRequest) (*v1.GetShipOrderStatusReply, error) {
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
	orderStatus, err := s.oc.GetShipOrderStatus(ctx, &biz.GetShipOrderStatusRequest{
		UserId:     userId,
		SubOrderId: req.SubOrderId,
	})
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "order does not belong to user")
	}
	log.Debugf("orderStatus: %+v", orderStatus)

	shippingAddressMap := map[string]any{
		"addressType":   orderStatus.ShippingAddress.AddressType,
		"city":          orderStatus.ShippingAddress.City,
		"contactPerson": orderStatus.ShippingAddress.ContactPerson,
		"contactPhone":  orderStatus.ShippingAddress.ContactPhone,
		"country":       orderStatus.ShippingAddress.Country,
		"state":         orderStatus.ShippingAddress.State,
		"streetAddress": orderStatus.ShippingAddress.StreetAddress,
		"zipCode":       orderStatus.ShippingAddress.ZipCode,
	}
	merchantAddress, err := structpb.NewStruct(shippingAddressMap)
	if err != nil {
		log.Errorf("failed to convert merchant address to struct: %v", err)
		return nil, status.Error(codes.Internal, "failed to convert merchant address to struct")
	}

	userAddressMap := map[string]any{
		"city":          orderStatus.ReceiverAddress.City,
		"country":       orderStatus.ReceiverAddress.Country,
		"createdAt":     orderStatus.ReceiverAddress.CreatedAt.Format(time.RFC3339),
		"email":         orderStatus.ReceiverAddress.Email,
		"id":            orderStatus.ReceiverAddress.ID,
		"state":         orderStatus.ReceiverAddress.State,
		"streetAddress": orderStatus.ReceiverAddress.StreetAddress,
		"updatedAt":     orderStatus.ReceiverAddress.UpdatedAt.Format(time.RFC3339),
		"userId":        orderStatus.ReceiverAddress.UserID,
		"zipCode":       orderStatus.ReceiverAddress.ZipCode,
	}
	userAddress, err := structpb.NewStruct(userAddressMap)
	if err != nil {
		// This is where the original error occurred
		log.Errorf("failed to convert user address to struct: %v", err)
		return nil, status.Error(codes.Internal, "failed to convert user address to struct")
	}
	log.Debugf("userAddress: %+v", userAddress)
	return &v1.GetShipOrderStatusReply{
		OrderId:         orderStatus.Id,
		SubOrderId:      orderStatus.SubOrderId,
		ShippingStatus:  pkg.MapShippingStatusToProto(orderStatus.ShippingStatus),
		ReceiverAddress: userAddress,
		ShippingAddress: merchantAddress,
		TrackingNumber:  orderStatus.TrackingNumber,
		Carrier:         orderStatus.Carrier,
	}, nil
}

// GetConsumerSubOrderDetail 根据用户id和子订单ID查询子订单详情
func (s *ConsumerOrderService) GetConsumerSubOrderDetail(ctx context.Context, req *v1.GetConsumerSubOrderDetailRequest) (*v1.Order, error) {
	userId, err := globalpkg.GetMetadataUesrID(ctx)
	order, err := s.oc.GetConsumerSubOrderDetail(ctx, &biz.GetConsumerSubOrderDetailRequest{
		UserId:     userId,
		SubOrderId: req.SubOrderId,
	})
	if err != nil {
		return nil, err
	}

	items := make([]*v1.OrderItem, 0, len(order.Items))
	for _, i := range order.Items {
		items = append(items, &v1.OrderItem{
			Item: &cartv1.CartItem{
				MerchantId: i.Item.MerchantId.String(),
				ProductId:  i.Item.ProductId.String(),
				Quantity:   i.Item.Quantity,
				Name:       i.Item.Name,
				Picture:    i.Item.Picture,
			},
			Cost: i.Cost,
		})
	}

	return &v1.Order{
		Items:      items,
		OrderId:    &order.OrderId,
		SubOrderId: &order.SubOrderID,
		Currency:   order.Currency,
		Address: &userv1.ConsumerAddress{
			City:          order.Address.City,
			State:         order.Address.State,
			Country:       order.Address.Country,
			ZipCode:       order.Address.ZipCode,
			StreetAddress: order.Address.StreetAddress,
		},
		Email:          order.Email,
		CreatedAt:      timestamppb.New(order.CreatedAt),
		PaymentStatus:  pkg.MapPaymentStatusToProto(order.PaymentStatus),
		ShippingStatus: pkg.MapShippingStatusToProto(order.ShippingStatus),
	}, nil
}
