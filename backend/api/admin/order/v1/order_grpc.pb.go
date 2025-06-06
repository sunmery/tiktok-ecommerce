// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: order/v1/order.proto

package adminorderv1

import (
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
	AdminOrder_GetAllOrders_FullMethodName = "/ecommerce.adminorder.v1.AdminOrder/GetAllOrders"
)

// AdminOrderClient is the client API for AdminOrder service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AdminOrderClient interface {
	// 查询全部订单列表
	GetAllOrders(ctx context.Context, in *GetAllOrdersReq, opts ...grpc.CallOption) (*AdminOrderReply, error)
}

type adminOrderClient struct {
	cc grpc.ClientConnInterface
}

func NewAdminOrderClient(cc grpc.ClientConnInterface) AdminOrderClient {
	return &adminOrderClient{cc}
}

func (c *adminOrderClient) GetAllOrders(ctx context.Context, in *GetAllOrdersReq, opts ...grpc.CallOption) (*AdminOrderReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AdminOrderReply)
	err := c.cc.Invoke(ctx, AdminOrder_GetAllOrders_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdminOrderServer is the server API for AdminOrder service.
// All implementations must embed UnimplementedAdminOrderServer
// for forward compatibility.
type AdminOrderServer interface {
	// 查询全部订单列表
	GetAllOrders(context.Context, *GetAllOrdersReq) (*AdminOrderReply, error)
	mustEmbedUnimplementedAdminOrderServer()
}

// UnimplementedAdminOrderServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAdminOrderServer struct{}

func (UnimplementedAdminOrderServer) GetAllOrders(context.Context, *GetAllOrdersReq) (*AdminOrderReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllOrders not implemented")
}
func (UnimplementedAdminOrderServer) mustEmbedUnimplementedAdminOrderServer() {}
func (UnimplementedAdminOrderServer) testEmbeddedByValue()                    {}

// UnsafeAdminOrderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AdminOrderServer will
// result in compilation errors.
type UnsafeAdminOrderServer interface {
	mustEmbedUnimplementedAdminOrderServer()
}

func RegisterAdminOrderServer(s grpc.ServiceRegistrar, srv AdminOrderServer) {
	// If the following call pancis, it indicates UnimplementedAdminOrderServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AdminOrder_ServiceDesc, srv)
}

func _AdminOrder_GetAllOrders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllOrdersReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminOrderServer).GetAllOrders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AdminOrder_GetAllOrders_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminOrderServer).GetAllOrders(ctx, req.(*GetAllOrdersReq))
	}
	return interceptor(ctx, in, info, handler)
}

// AdminOrder_ServiceDesc is the grpc.ServiceDesc for AdminOrder service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AdminOrder_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ecommerce.adminorder.v1.AdminOrder",
	HandlerType: (*AdminOrderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAllOrders",
			Handler:    _AdminOrder_GetAllOrders_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "order/v1/order.proto",
}
