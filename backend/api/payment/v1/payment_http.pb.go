// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.4
// - protoc             v5.29.3
// source: v1/payment.proto

package paymentv1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationPaymentServiceCreatePayment = "/payment.v1.PaymentService/CreatePayment"
const OperationPaymentServiceGetPaymentStatus = "/payment.v1.PaymentService/GetPaymentStatus"
const OperationPaymentServiceHandlePaymentCallback = "/payment.v1.PaymentService/HandlePaymentCallback"
const OperationPaymentServiceHandlePaymentNotify = "/payment.v1.PaymentService/HandlePaymentNotify"

type PaymentServiceHTTPServer interface {
	// CreatePayment 创建支付订单
	CreatePayment(context.Context, *CreatePaymentRequest) (*CreatePaymentResponse, error)
	// GetPaymentStatus 查询支付状态
	GetPaymentStatus(context.Context, *GetPaymentStatusRequest) (*GetPaymentStatusResponse, error)
	// HandlePaymentCallback 支付成功后的回调处理
	HandlePaymentCallback(context.Context, *HandlePaymentCallbackRequest) (*HandlePaymentCallbackResponse, error)
	// HandlePaymentNotify 处理支付回调通知
	//  rpc HandlePaymentNotify (HandlePaymentNotifyRequest) returns (HandlePaymentNotifyResponse) {
	HandlePaymentNotify(context.Context, *UrlValues) (*HandlePaymentNotifyResponse, error)
}

func RegisterPaymentServiceHTTPServer(s *http.Server, srv PaymentServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/v1/payments", _PaymentService_CreatePayment0_HTTP_Handler(srv))
	r.GET("/v1/payments/{payment_id}/status", _PaymentService_GetPaymentStatus0_HTTP_Handler(srv))
	r.POST("/v1/payments/notify", _PaymentService_HandlePaymentNotify0_HTTP_Handler(srv))
	r.GET("/v1/payments/callback", _PaymentService_HandlePaymentCallback0_HTTP_Handler(srv))
}

func _PaymentService_CreatePayment0_HTTP_Handler(srv PaymentServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreatePaymentRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPaymentServiceCreatePayment)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreatePayment(ctx, req.(*CreatePaymentRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreatePaymentResponse)
		return ctx.Result(200, reply)
	}
}

func _PaymentService_GetPaymentStatus0_HTTP_Handler(srv PaymentServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetPaymentStatusRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPaymentServiceGetPaymentStatus)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetPaymentStatus(ctx, req.(*GetPaymentStatusRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetPaymentStatusResponse)
		return ctx.Result(200, reply)
	}
}

func _PaymentService_HandlePaymentNotify0_HTTP_Handler(srv PaymentServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UrlValues
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPaymentServiceHandlePaymentNotify)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.HandlePaymentNotify(ctx, req.(*UrlValues))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*HandlePaymentNotifyResponse)
		return ctx.Result(200, reply)
	}
}

func _PaymentService_HandlePaymentCallback0_HTTP_Handler(srv PaymentServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in HandlePaymentCallbackRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPaymentServiceHandlePaymentCallback)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.HandlePaymentCallback(ctx, req.(*HandlePaymentCallbackRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*HandlePaymentCallbackResponse)
		return ctx.Result(200, reply)
	}
}

type PaymentServiceHTTPClient interface {
	CreatePayment(ctx context.Context, req *CreatePaymentRequest, opts ...http.CallOption) (rsp *CreatePaymentResponse, err error)
	GetPaymentStatus(ctx context.Context, req *GetPaymentStatusRequest, opts ...http.CallOption) (rsp *GetPaymentStatusResponse, err error)
	HandlePaymentCallback(ctx context.Context, req *HandlePaymentCallbackRequest, opts ...http.CallOption) (rsp *HandlePaymentCallbackResponse, err error)
	HandlePaymentNotify(ctx context.Context, req *UrlValues, opts ...http.CallOption) (rsp *HandlePaymentNotifyResponse, err error)
}

type PaymentServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewPaymentServiceHTTPClient(client *http.Client) PaymentServiceHTTPClient {
	return &PaymentServiceHTTPClientImpl{client}
}

func (c *PaymentServiceHTTPClientImpl) CreatePayment(ctx context.Context, in *CreatePaymentRequest, opts ...http.CallOption) (*CreatePaymentResponse, error) {
	var out CreatePaymentResponse
	pattern := "/v1/payments"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationPaymentServiceCreatePayment))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *PaymentServiceHTTPClientImpl) GetPaymentStatus(ctx context.Context, in *GetPaymentStatusRequest, opts ...http.CallOption) (*GetPaymentStatusResponse, error) {
	var out GetPaymentStatusResponse
	pattern := "/v1/payments/{payment_id}/status"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationPaymentServiceGetPaymentStatus))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *PaymentServiceHTTPClientImpl) HandlePaymentCallback(ctx context.Context, in *HandlePaymentCallbackRequest, opts ...http.CallOption) (*HandlePaymentCallbackResponse, error) {
	var out HandlePaymentCallbackResponse
	pattern := "/v1/payments/callback"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationPaymentServiceHandlePaymentCallback))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *PaymentServiceHTTPClientImpl) HandlePaymentNotify(ctx context.Context, in *UrlValues, opts ...http.CallOption) (*HandlePaymentNotifyResponse, error) {
	var out HandlePaymentNotifyResponse
	pattern := "/v1/payments/notify"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationPaymentServiceHandlePaymentNotify))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
