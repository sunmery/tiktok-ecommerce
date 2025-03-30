// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: v1/user.proto

package userv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	UserService_GetUserProfile_FullMethodName   = "/ecommerce.user.v1.UserService/GetUserProfile"
	UserService_GetUsers_FullMethodName         = "/ecommerce.user.v1.UserService/GetUsers"
	UserService_DeleteUser_FullMethodName       = "/ecommerce.user.v1.UserService/DeleteUser"
	UserService_UpdateUser_FullMethodName       = "/ecommerce.user.v1.UserService/UpdateUser"
	UserService_CreateAddresses_FullMethodName  = "/ecommerce.user.v1.UserService/CreateAddresses"
	UserService_UpdateAddresses_FullMethodName  = "/ecommerce.user.v1.UserService/UpdateAddresses"
	UserService_DeleteAddresses_FullMethodName  = "/ecommerce.user.v1.UserService/DeleteAddresses"
	UserService_GetAddress_FullMethodName       = "/ecommerce.user.v1.UserService/GetAddress"
	UserService_CreateCreditCard_FullMethodName = "/ecommerce.user.v1.UserService/CreateCreditCard"
	UserService_GetAddresses_FullMethodName     = "/ecommerce.user.v1.UserService/GetAddresses"
	UserService_ListCreditCards_FullMethodName  = "/ecommerce.user.v1.UserService/ListCreditCards"
	UserService_GetCreditCard_FullMethodName    = "/ecommerce.user.v1.UserService/GetCreditCard"
	UserService_DeleteCreditCard_FullMethodName = "/ecommerce.user.v1.UserService/DeleteCreditCard"
)

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// 用户服务接口定义
type UserServiceClient interface {
	// 获取用户个人资料
	GetUserProfile(ctx context.Context, in *GetProfileRequest, opts ...grpc.CallOption) (*GetProfileResponse, error)
	// 获取全部用户信息
	GetUsers(ctx context.Context, in *GetUsersRequest, opts ...grpc.CallOption) (*GetUsersResponse, error)
	// 删除用户
	DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*DeleteUserResponse, error)
	// 更新用户信息
	UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*UpdateUserResponse, error)
	// 创建用户地址
	CreateAddresses(ctx context.Context, in *Address, opts ...grpc.CallOption) (*Address, error)
	// 更新用户地址
	UpdateAddresses(ctx context.Context, in *Address, opts ...grpc.CallOption) (*Address, error)
	// 删除用户地址
	DeleteAddresses(ctx context.Context, in *DeleteAddressesRequest, opts ...grpc.CallOption) (*DeleteAddressesReply, error)
	// 根据 ID获取用户地址
	GetAddress(ctx context.Context, in *GetAddressRequest, opts ...grpc.CallOption) (*Address, error)
	// 创建用户的信用卡信息
	CreateCreditCard(ctx context.Context, in *CreditCard, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// 获取用户地址列表
	GetAddresses(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetAddressesReply, error)
	// 列出用户的信用卡信息
	ListCreditCards(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*CreditCards, error)
	// 根据ID搜索用户的信用卡信息
	GetCreditCard(ctx context.Context, in *GetCreditCardRequest, opts ...grpc.CallOption) (*CreditCard, error)
	// 删除用户的信用卡信息
	DeleteCreditCard(ctx context.Context, in *DeleteCreditCardsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) GetUserProfile(ctx context.Context, in *GetProfileRequest, opts ...grpc.CallOption) (*GetProfileResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetProfileResponse)
	err := c.cc.Invoke(ctx, UserService_GetUserProfile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetUsers(ctx context.Context, in *GetUsersRequest, opts ...grpc.CallOption) (*GetUsersResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUsersResponse)
	err := c.cc.Invoke(ctx, UserService_GetUsers_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*DeleteUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteUserResponse)
	err := c.cc.Invoke(ctx, UserService_DeleteUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*UpdateUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateUserResponse)
	err := c.cc.Invoke(ctx, UserService_UpdateUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) CreateAddresses(ctx context.Context, in *Address, opts ...grpc.CallOption) (*Address, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Address)
	err := c.cc.Invoke(ctx, UserService_CreateAddresses_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UpdateAddresses(ctx context.Context, in *Address, opts ...grpc.CallOption) (*Address, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Address)
	err := c.cc.Invoke(ctx, UserService_UpdateAddresses_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) DeleteAddresses(ctx context.Context, in *DeleteAddressesRequest, opts ...grpc.CallOption) (*DeleteAddressesReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteAddressesReply)
	err := c.cc.Invoke(ctx, UserService_DeleteAddresses_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetAddress(ctx context.Context, in *GetAddressRequest, opts ...grpc.CallOption) (*Address, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Address)
	err := c.cc.Invoke(ctx, UserService_GetAddress_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) CreateCreditCard(ctx context.Context, in *CreditCard, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, UserService_CreateCreditCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetAddresses(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetAddressesReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAddressesReply)
	err := c.cc.Invoke(ctx, UserService_GetAddresses_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) ListCreditCards(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*CreditCards, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreditCards)
	err := c.cc.Invoke(ctx, UserService_ListCreditCards_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetCreditCard(ctx context.Context, in *GetCreditCardRequest, opts ...grpc.CallOption) (*CreditCard, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreditCard)
	err := c.cc.Invoke(ctx, UserService_GetCreditCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) DeleteCreditCard(ctx context.Context, in *DeleteCreditCardsRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, UserService_DeleteCreditCard_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility.
//
// 用户服务接口定义
type UserServiceServer interface {
	// 获取用户个人资料
	GetUserProfile(context.Context, *GetProfileRequest) (*GetProfileResponse, error)
	// 获取全部用户信息
	GetUsers(context.Context, *GetUsersRequest) (*GetUsersResponse, error)
	// 删除用户
	DeleteUser(context.Context, *DeleteUserRequest) (*DeleteUserResponse, error)
	// 更新用户信息
	UpdateUser(context.Context, *UpdateUserRequest) (*UpdateUserResponse, error)
	// 创建用户地址
	CreateAddresses(context.Context, *Address) (*Address, error)
	// 更新用户地址
	UpdateAddresses(context.Context, *Address) (*Address, error)
	// 删除用户地址
	DeleteAddresses(context.Context, *DeleteAddressesRequest) (*DeleteAddressesReply, error)
	// 根据 ID获取用户地址
	GetAddress(context.Context, *GetAddressRequest) (*Address, error)
	// 创建用户的信用卡信息
	CreateCreditCard(context.Context, *CreditCard) (*emptypb.Empty, error)
	// 获取用户地址列表
	GetAddresses(context.Context, *emptypb.Empty) (*GetAddressesReply, error)
	// 列出用户的信用卡信息
	ListCreditCards(context.Context, *emptypb.Empty) (*CreditCards, error)
	// 根据ID搜索用户的信用卡信息
	GetCreditCard(context.Context, *GetCreditCardRequest) (*CreditCard, error)
	// 删除用户的信用卡信息
	DeleteCreditCard(context.Context, *DeleteCreditCardsRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedUserServiceServer struct{}

func (UnimplementedUserServiceServer) GetUserProfile(context.Context, *GetProfileRequest) (*GetProfileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserProfile not implemented")
}
func (UnimplementedUserServiceServer) GetUsers(context.Context, *GetUsersRequest) (*GetUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUsers not implemented")
}
func (UnimplementedUserServiceServer) DeleteUser(context.Context, *DeleteUserRequest) (*DeleteUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}
func (UnimplementedUserServiceServer) UpdateUser(context.Context, *UpdateUserRequest) (*UpdateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}
func (UnimplementedUserServiceServer) CreateAddresses(context.Context, *Address) (*Address, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAddresses not implemented")
}
func (UnimplementedUserServiceServer) UpdateAddresses(context.Context, *Address) (*Address, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAddresses not implemented")
}
func (UnimplementedUserServiceServer) DeleteAddresses(context.Context, *DeleteAddressesRequest) (*DeleteAddressesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAddresses not implemented")
}
func (UnimplementedUserServiceServer) GetAddress(context.Context, *GetAddressRequest) (*Address, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAddress not implemented")
}
func (UnimplementedUserServiceServer) CreateCreditCard(context.Context, *CreditCard) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCreditCard not implemented")
}
func (UnimplementedUserServiceServer) GetAddresses(context.Context, *emptypb.Empty) (*GetAddressesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAddresses not implemented")
}
func (UnimplementedUserServiceServer) ListCreditCards(context.Context, *emptypb.Empty) (*CreditCards, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCreditCards not implemented")
}
func (UnimplementedUserServiceServer) GetCreditCard(context.Context, *GetCreditCardRequest) (*CreditCard, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCreditCard not implemented")
}
func (UnimplementedUserServiceServer) DeleteCreditCard(context.Context, *DeleteCreditCardsRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCreditCard not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}
func (UnimplementedUserServiceServer) testEmbeddedByValue()                     {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	// If the following call pancis, it indicates UnimplementedUserServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_GetUserProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProfileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUserProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetUserProfile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUserProfile(ctx, req.(*GetProfileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUsers(ctx, req.(*GetUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_DeleteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).DeleteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_DeleteUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).DeleteUser(ctx, req.(*DeleteUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_UpdateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UpdateUser(ctx, req.(*UpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_CreateAddresses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Address)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CreateAddresses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_CreateAddresses_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).CreateAddresses(ctx, req.(*Address))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UpdateAddresses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Address)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UpdateAddresses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_UpdateAddresses_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UpdateAddresses(ctx, req.(*Address))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_DeleteAddresses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAddressesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).DeleteAddresses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_DeleteAddresses_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).DeleteAddresses(ctx, req.(*DeleteAddressesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetAddress_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetAddress(ctx, req.(*GetAddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_CreateCreditCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreditCard)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CreateCreditCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_CreateCreditCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).CreateCreditCard(ctx, req.(*CreditCard))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetAddresses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetAddresses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetAddresses_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetAddresses(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_ListCreditCards_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).ListCreditCards(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_ListCreditCards_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).ListCreditCards(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetCreditCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCreditCardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetCreditCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetCreditCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetCreditCard(ctx, req.(*GetCreditCardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_DeleteCreditCard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCreditCardsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).DeleteCreditCard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_DeleteCreditCard_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).DeleteCreditCard(ctx, req.(*DeleteCreditCardsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ecommerce.user.v1.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserProfile",
			Handler:    _UserService_GetUserProfile_Handler,
		},
		{
			MethodName: "GetUsers",
			Handler:    _UserService_GetUsers_Handler,
		},
		{
			MethodName: "DeleteUser",
			Handler:    _UserService_DeleteUser_Handler,
		},
		{
			MethodName: "UpdateUser",
			Handler:    _UserService_UpdateUser_Handler,
		},
		{
			MethodName: "CreateAddresses",
			Handler:    _UserService_CreateAddresses_Handler,
		},
		{
			MethodName: "UpdateAddresses",
			Handler:    _UserService_UpdateAddresses_Handler,
		},
		{
			MethodName: "DeleteAddresses",
			Handler:    _UserService_DeleteAddresses_Handler,
		},
		{
			MethodName: "GetAddress",
			Handler:    _UserService_GetAddress_Handler,
		},
		{
			MethodName: "CreateCreditCard",
			Handler:    _UserService_CreateCreditCard_Handler,
		},
		{
			MethodName: "GetAddresses",
			Handler:    _UserService_GetAddresses_Handler,
		},
		{
			MethodName: "ListCreditCards",
			Handler:    _UserService_ListCreditCards_Handler,
		},
		{
			MethodName: "GetCreditCard",
			Handler:    _UserService_GetCreditCard_Handler,
		},
		{
			MethodName: "DeleteCreditCard",
			Handler:    _UserService_DeleteCreditCard_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/user.proto",
}
