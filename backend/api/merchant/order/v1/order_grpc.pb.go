// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: order/v1/order.proto

package orderv1

import (
	v1 "backend/api/order/v1"
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Order_GetMerchantOrders_FullMethodName = "/ecommerce.merchant.v1.Order/GetMerchantOrders"
)

// OrderClient is the client API for Order service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrderClient interface {
	// 查询商家订单列表(商家侧)
	GetMerchantOrders(ctx context.Context, in *GetMerchantOrdersReq, opts ...grpc.CallOption) (*v1.Orders, error)
}

type orderClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderClient(cc grpc.ClientConnInterface) OrderClient {
	return &orderClient{cc}
}

func (c *orderClient) GetMerchantOrders(ctx context.Context, in *GetMerchantOrdersReq, opts ...grpc.CallOption) (*v1.Orders, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(v1.Orders)
	err := c.cc.Invoke(ctx, Order_GetMerchantOrders_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrderServer is the server API for Order service.
// All implementations must embed UnimplementedOrderServer
// for forward compatibility.
type OrderServer interface {
	// 查询商家订单列表(商家侧)
	GetMerchantOrders(context.Context, *GetMerchantOrdersReq) (*v1.Orders, error)
	mustEmbedUnimplementedOrderServer()
}

// UnimplementedOrderServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedOrderServer struct{}

func (UnimplementedOrderServer) GetMerchantOrders(context.Context, *GetMerchantOrdersReq) (*v1.Orders, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMerchantOrders not implemented")
}
func (UnimplementedOrderServer) mustEmbedUnimplementedOrderServer() {}
func (UnimplementedOrderServer) testEmbeddedByValue()               {}

// UnsafeOrderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrderServer will
// result in compilation errors.
type UnsafeOrderServer interface {
	mustEmbedUnimplementedOrderServer()
}

func RegisterOrderServer(s grpc.ServiceRegistrar, srv OrderServer) {
	// If the following call pancis, it indicates UnimplementedOrderServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Order_ServiceDesc, srv)
}

func _Order_GetMerchantOrders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMerchantOrdersReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServer).GetMerchantOrders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Order_GetMerchantOrders_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServer).GetMerchantOrders(ctx, req.(*GetMerchantOrdersReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Order_ServiceDesc is the grpc.ServiceDesc for Order service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Order_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ecommerce.merchant.v1.Order",
	HandlerType: (*OrderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetMerchantOrders",
			Handler:    _Order_GetMerchantOrders_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "order/v1/order.proto",
}
