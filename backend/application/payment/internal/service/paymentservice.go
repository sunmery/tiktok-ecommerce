package service

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/uuid"

	"github.com/go-kratos/kratos/v2/transport"

	"github.com/go-kratos/kratos/v2/transport/http"

	"backend/application/payment/internal/biz"
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
	consumerId, err := uuid.Parse(req.ConsumerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "无效的用户ID")
	}

	log.Debugf("consumerId%v", consumerId)

	// merchantId, err := uuid.Parse(req.MerchantId)
	// if err != nil {
	// 	return nil, status.Error(codes.InvalidArgument, "无效的商家ID")
	// }

	// 从上下文或请求获取订单ID
	var orderID int64
	if req.OrderId != 0 {
		orderID = req.OrderId
	} else {
		return nil, status.Error(codes.InvalidArgument, "订单ID不能为空")
	}

	// 创建支付请求
	createReq := &biz.CreatePaymentReq{
		OrderID:    orderID,
		ConsumerID: consumerId,
		// MerchantID:      merchantId,
		Amount:          req.Amount,
		Currency:        req.Currency,
		Subject:         req.Subject,
		ReturnURL:       req.ReturnUrl,
		FreezeId:        req.FreezeId,
		ConsumerVersion: req.ConsumerVersion,
		MerchanVersion:  req.MerchantVersion,
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
func (s *PaymentService) HandlePaymentNotify(ctx context.Context, req *pb.UrlValues) (*pb.HandlePaymentNotifyResponse, error) {
	s.log.WithContext(ctx).Infof("service HandlePaymentNotify: %v", req)

	values := make(url.Values)
	// 安全地获取HTTP上下文
	httpCtx, ok := ctx.(http.Context)
	if !ok {
		s.log.WithContext(ctx).Errorf("无法获取HTTP上下文，使用请求参数继续处理")
		// 即使没有HTTP上下文，也继续处理，使用请求中的参数
	} else {
		// 确保解析表单数据
		if err := httpCtx.Request().ParseForm(); err != nil {
			s.log.WithContext(ctx).Errorf("解析表单数据失败: %v", err)
		}

		// 安全地获取表单数据
		form := httpCtx.Form()
		if form != nil {
			for k, v := range form {
				fmt.Printf("k: %v, v: %v\n", k, v)
				values[k] = v
			}
		}
	}
	// ProtoToUrlValues(values)
	s.log.WithContext(ctx).Infof("service HandlePaymentNotify values: %v", values)
	// 转换请求
	// notifyReq := &biz.PaymentNotifyReq{
	// 	AppID:       req.AppId,
	// 	AuthAppId:   req.AuthAppId,
	// 	TradeNo:     req.TradeNo,
	// 	Charset:     req.Charset,
	// 	Method:      req.Method,
	// 	Sign:        req.Sign,
	// 	SignType:    req.SignType,
	// 	OutTradeNo:  req.OutTradeNo,
	// 	TotalAmount: req.TotalAmount,
	// 	SellerId:    req.SellerId,
	// 	Params:      values,
	// }

	// 调用业务逻辑
	resp, err := s.uc.HandlePaymentNotify(ctx, values)
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

func ProtoToUrlValues(protoData *pb.UrlValues) url.Values {
	values := url.Values{}
	for _, pair := range protoData.Pairs {
		for _, v := range pair.Values {
			values.Add(pair.Key, v)
		}
	}
	return values
}

func copyValues(dst, src url.Values) {
	for k, vs := range src {
		dst[k] = append(dst[k], vs...)
	}
}

// HandlePaymentCallback 处理支付回调
func (s *PaymentService) HandlePaymentCallback(ctx context.Context, req *pb.HandlePaymentCallbackRequest) (*pb.HandlePaymentCallbackResponse, error) {
	s.log.WithContext(ctx).Infof("service HandlePaymentCallback: %v", req)

	// 将请求参数转换为url.Values格式
	values := make(url.Values)

	// 安全地获取HTTP上下文
	tr, ok := transport.FromServerContext(ctx)
	if !ok {
		s.log.WithContext(ctx).Errorf("无法获取HTTP上下文，使用请求参数继续处理")
	}
	httpCtx, ok := tr.(http.Context)
	if !ok {
		s.log.WithContext(ctx).Errorf("无法获取HTTP上下文，使用请求参数继续处理")
		// 即使没有HTTP上下文，也继续处理，使用请求中的参数
	} else {
		// 确保解析表单数据
		if err := httpCtx.Request().ParseForm(); err != nil {
			s.log.WithContext(ctx).Errorf("解析表单数据失败: %v", err)
		}

		// 安全地获取表单数据
		form := httpCtx.Form()
		for k, v := range form {
			fmt.Printf("k: %v, v: %v\n", k, v)
			values[k] = v
		}

	}
	s.log.WithContext(ctx).Infof("service values: %v", values)

	// 转换请求
	callbackReq := &biz.PaymentCallbackReq{
		Params:      req.Params,
		OutTradeNo:  req.OutTradeNo,
		TradeNo:     req.TradeNo,
		TotalAmount: req.TotalAmount,
		Subject:     req.Subject,
		TradeStatus: req.TradeStatus,
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
