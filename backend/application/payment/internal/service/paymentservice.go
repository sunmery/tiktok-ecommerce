package service

import (
	"context"

	"backend/application/payment/internal/biz"
	"backend/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "backend/api/payment/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// PaymentService 支付服务实现
type PaymentService struct {
	pb.UnimplementedPaymentServiceServer
	uc  *biz.PaymentUsecase
	log *log.Helper
}

func NewPaymentService(uc *biz.PaymentUsecase, logger log.Logger) *PaymentService {
	return &PaymentService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// CreatePayment 创建支付订单
func (s *PaymentService) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	s.log.WithContext(ctx).Infof("CreatePayment: %v", req)

	uid, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "无效的用户ID")
	}

	// 从上下文或请求获取订单ID
	var orderID int64
	if req.OrderId != 0 {
		orderID = req.OrderId
	} else {
		return nil, status.Error(codes.InvalidArgument, "订单ID不能为空")
	}

	// 创建支付请求
	createReq := &biz.CreatePaymentReq{
		OrderID:   orderID,
		UserID:    uid,
		Amount:    req.Amount,
		Currency:  req.Currency,
		Subject:   req.Subject,
		ReturnURL: req.ReturnUrl,
	}

	// 调用业务逻辑
	result, err := s.uc.CreatePayment(ctx, createReq)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to create payment: %v", err)
		return nil, err
	}

	// 返回创建结果
	return &pb.CreatePaymentResponse{
		PaymentId: result.Payment.ID,
		PayUrl:    result.Payment.PayURL,
	}, nil
}

// GetPaymentStatus 查询支付状态
func (s *PaymentService) GetPaymentStatus(ctx context.Context, req *pb.GetPaymentStatusRequest) (*pb.GetPaymentStatusResponse, error) {
	s.log.WithContext(ctx).Infof("GetPaymentStatus: %v", req)

	// 转换请求
	getReq := &biz.GetPaymentStatusReq{
		PaymentID: req.PaymentId,
	}

	// 调用业务逻辑
	resp, err := s.uc.GetPaymentStatus(ctx, getReq)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to get payment status: %v", err)
		return nil, err
	}

	// 转换支付状态
	var payStatus pb.PaymentStatus
	switch resp.Payment.Status {
	case biz.PaymentStatusPending:
		payStatus = pb.PaymentStatus_PAYMENT_STATUS_PENDING
	case biz.PaymentStatusProcessing:
		payStatus = pb.PaymentStatus_PAYMENT_STATUS_PROCESSING
	case biz.PaymentStatusSuccess:
		payStatus = pb.PaymentStatus_PAYMENT_STATUS_SUCCESS
	case biz.PaymentStatusFailed:
		payStatus = pb.PaymentStatus_PAYMENT_STATUS_FAILED
	case biz.PaymentStatusClosed:
		payStatus = pb.PaymentStatus_PAYMENT_STATUS_CLOSED
	default:
		payStatus = pb.PaymentStatus_PAYMENT_STATUS_UNKNOWN
	}

	// 返回结果
	return &pb.GetPaymentStatusResponse{
		PaymentId: resp.Payment.ID,
		OrderId:   resp.Payment.OrderID,
		Status:    payStatus,
		TradeNo:   resp.Payment.TradeNo,
	}, nil
}

// HandlePaymentNotify 处理支付通知
func (s *PaymentService) HandlePaymentNotify(ctx context.Context, req *pb.HandlePaymentNotifyRequest) (*pb.HandlePaymentNotifyResponse, error) {
	s.log.WithContext(ctx).Infof("HandlePaymentNotify: %v", req)

	// 转换请求
	notifyReq := &biz.PaymentNotifyReq{
		AppID:       req.AppId,
		TradeNo:     req.TradeNo,
		OutTradeNo:  req.OutTradeNo,
		TotalAmount: req.TotalAmount,
		Subject:     req.Subject,
		TradeStatus: req.TradeStatus,
		GmtPayment:  req.GmtPayment,
		GmtCreate:   req.GmtCreate,
		Sign:        req.Sign,
		SignType:    req.SignType,
		Params:      req.Params,
	}

	// 调用业务逻辑
	resp, err := s.uc.HandlePaymentNotify(ctx, notifyReq)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to handle payment notify: %v", err)
		return nil, err
	}

	// 返回结果
	return &pb.HandlePaymentNotifyResponse{
		Success: resp.Success,
		Message: resp.Message,
	}, nil
}

// HandlePaymentCallback 处理支付回调
func (s *PaymentService) HandlePaymentCallback(ctx context.Context, req *pb.HandlePaymentCallbackRequest) (*pb.HandlePaymentCallbackResponse, error) {
	s.log.WithContext(ctx).Infof("HandlePaymentCallback: %v", req)

	// 转换请求
	callbackReq := &biz.PaymentCallbackReq{
		OutTradeNo:  req.OutTradeNo,
		TradeNo:     req.TradeNo,
		TotalAmount: req.TotalAmount,
		Params:      req.Params,
	}

	// 调用业务逻辑
	resp, err := s.uc.HandlePaymentCallback(ctx, callbackReq)
	if err != nil {
		s.log.WithContext(ctx).Errorf("Failed to handle payment callback: %v", err)
		return nil, err
	}

	// 返回结果
	return &pb.HandlePaymentCallbackResponse{
		Success: resp.Success,
		Message: resp.Message,
	}, nil
}

// func (s *PaymentService) GetPayment(ctx context.Context, req *pb.GetPaymentReq) (*pb.PaymentResp, error) {
// 	paymentId, err := pkg.GetMetadataUesrID(ctx)
// 	if err != nil {
// 		return nil, status.Error(codes.InvalidArgument, "invalid payment ID")
// 	}
// 	payment, err := s.uc.GetPayment(ctx, paymentId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &pb.PaymentResp{
// 		PaymentId:  payment.PaymentId,
// 		Status:     payment.Status,
// 		PaymentUrl: payment.PaymentUrl,
// 		CreatedAt:  timestamppb.New(payment.CreatedAt),
// 	}, nil
// }
