package service

import (
	"context"

	"backend/constants"

	"backend/application/merchant/internal/pkg"

	globalPkg "backend/pkg"

	orderv1 "backend/api/merchant/order/v1"

	"github.com/google/uuid"

	"github.com/go-kratos/kratos/v2/log"

	cartv1 "backend/api/cart/v1"
	v1 "backend/api/order/v1"
	userv1 "backend/api/user/v1"

	"google.golang.org/protobuf/types/known/timestamppb"

	"backend/application/merchant/internal/biz"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderService struct {
	orderv1.UnimplementedOrderServer
	oc *biz.OrderUsecase
}

func NewOrderService(oc *biz.OrderUsecase) *OrderService {
	return &OrderService{oc: oc}
}

// GetMerchantOrders 获取商家订单列表
func (s *OrderService) GetMerchantOrders(ctx context.Context, req *orderv1.GetMerchantOrdersReq) (*v1.Orders, error) {
	// 从网关获取用户ID
	userId, err := globalPkg.GetMetadataUesrID(ctx)
	if err != nil {
		log.Errorf("获取用户ID失败: %v", err)
		return nil, status.Error(codes.Unauthenticated, "获取用户ID失败")
	}

	// 使用请求中的merchant_id覆盖，如果有指定的话
	if req.MerchantId != "" {
		// 此处可以添加权限检查，确保当前用户有权查看指定商家的订单
		log.Infof("使用指定的商家ID: %s", req.MerchantId)
		parsedId, parseErr := uuid.Parse(req.MerchantId)
		if parseErr != nil {
			log.Errorf("商家ID格式无效: %v", parseErr)
			return nil, status.Error(codes.InvalidArgument, "商家ID格式无效")
		}
		userId = parsedId
	}

	// 调用业务层获取订单列表
	resp, err := s.oc.GetMerchantOrders(ctx, &biz.GetMerchantOrdersReq{
		UserID:   userId,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		log.Errorf("获取商家订单失败: %v", err)
		return nil, status.Errorf(codes.Internal, "获取商家订单失败: %v", err)
	}

	// 检查是否有订单
	if len(resp.Orders) == 0 {
		log.Infof("商家 %s 没有订单记录", userId)
		return &v1.Orders{Orders: []*v1.Order{}}, nil
	}

	// 按照商家订单分组
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
		var shippingStatus constants.ShippingStatus
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
			shippingStatus = subOrder.ShippingStatus
		}

		// 转换时间戳
		createdAt := timestamppb.New(firstSubOrder.CreatedAt)

		// 解析支付状态和运输状态
		paymentStatus := pkg.MapPaymentStatusToProto(string(firstSubOrder.PaymentStatus))
		ShippingStatus := pkg.MapShippingStatusToProto(shippingStatus)
		// 创建地址信息 (在真实场景中需要从订单数据中获取)
		address := &userv1.Address{
			StreetAddress: "未提供地址信息", // TODO 这里应该从订单数据中获取实际地址
			City:          "",
			State:         "",
			Country:       "",
			ZipCode:       "",
		}

		// 添加订单到响应列表
		orders = append(orders, &v1.Order{
			Items:          orderItems,
			OrderId:        firstSubOrder.ID, // 注意: 确保ID类型转换正确
			UserId:         firstSubOrder.MerchantID.String(),
			Currency:       firstSubOrder.Currency,
			Address:        address,
			Email:          "未提供邮箱", // TODO 这里应该从订单数据中获取实际邮箱
			CreatedAt:      createdAt,
			PaymentStatus:  paymentStatus,
			ShippingStatus: ShippingStatus,
		})
	}
	return &v1.Orders{Orders: orders}, nil
}
