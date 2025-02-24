package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"

	"google.golang.org/protobuf/types/known/timestamppb"

	"backend/application/payment/internal/biz"

	pb "backend/api/payment/v1"
)

type PaymentServiceService struct {
	pb.UnimplementedPaymentServiceServer

	uc *biz.PaymentUsecase
}

func NewPaymentServiceService(uc *biz.PaymentUsecase) *PaymentServiceService {
	return &PaymentServiceService{uc: uc}
}

func (s *PaymentServiceService) CreatePayment(ctx context.Context, req *pb.CreatePaymentReq) (*pb.PaymentResp, error) {
	orderId, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, fmt.Errorf("invalid order id")
	}
	result, err := s.uc.CreatePayment(ctx, &biz.CreatePaymentReq{
		OrderId:  orderId,
		Currency: req.Currency,
		Amount:   req.Amount,
		// PaymentMethod: req.PaymentMethod,
		// Metadata:      req.Metadata,
	})
	if err != nil {
		return nil, err
	}
	return &pb.PaymentResp{
		PaymentId:  result.PaymentId,
		Status:     result.Status,
		PaymentUrl: result.PaymentUrl,
		CreatedAt:  timestamppb.New(result.CreatedAt),
	}, nil
}

func (s *PaymentServiceService) PaymentNotify(ctx context.Context, req *pb.PaymentNotifyReq) (*pb.PaymentNotifyResp, error) {
	requestForm := make(map[string][]string)
	for k, v := range req.Values {
		requestForm[k] = v.Values
	}
	result, err := s.uc.PaymentNotify(ctx, &biz.PaymentNotifyReq{
		Values: requestForm,
	})
	if err != nil {
		return nil, err
	}
	fmt.Printf("result: %v\n", result)
	return &pb.PaymentNotifyResp{
		Code: strconv.Itoa(result.Code),
		Msg:  result.Msg,
	}, nil
}

func (s *PaymentServiceService) ProcessPaymentCallback(ctx context.Context, req *pb.PaymentCallbackReq) (*pb.PaymentCallbackResp, error) {
	requestForm := make(map[string][]string)
	for k, v := range req.RequestForm {
		requestForm[k] = v.Values
	}

	callback, err := s.uc.ProcessPaymentCallback(ctx, &biz.PaymentCallbackReq{
		PaymentId:       req.PaymentId,
		Status:          biz.PaymentStatus(req.Status),
		GatewayResponse: req.GatewayResponse,
		ProcessedAt:     req.ProcessedAt.AsTime(),
		RequestForm:     requestForm,
	})
	fmt.Printf("callback: %v\n", callback)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *PaymentServiceService) GetPayment(ctx context.Context, req *pb.GetPaymentReq) (*pb.PaymentResp, error) {
	paymentId, err := uuid.Parse(req.PaymentId)
	if err != nil {
		return nil, err
	}
	payment, err := s.uc.GetPayment(ctx, paymentId)
	if err != nil {
		return nil, err
	}
	return &pb.PaymentResp{
		PaymentId:  payment.PaymentId,
		Status:     payment.Status,
		PaymentUrl: payment.PaymentUrl,
		CreatedAt:  timestamppb.New(payment.CreatedAt),
	}, nil
}
