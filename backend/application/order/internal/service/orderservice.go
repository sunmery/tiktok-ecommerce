package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"backend/application/order/internal/biz"
	globalpkg "backend/pkg"

	v1 "backend/api/order/v1"

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

func (s *OrderServiceService) MarkOrderPaid(ctx context.Context, req *v1.MarkOrderPaidReq) (*v1.MarkOrderPaidResp, error) {
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
	orderPaid, markOrderPaidErr := s.uc.MarkOrderPaid(ctx, &biz.MarkOrderPaidReq{
		UserId:  userId,
		OrderId: req.OrderId,
	})
	if markOrderPaidErr != nil {
		return nil, markOrderPaidErr
	}

	log.Debugf("orderPaid: %v", orderPaid)

	return &v1.MarkOrderPaidResp{}, nil
}
