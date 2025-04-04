// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: v1/category.proto

// 定义包名，用于区分不同的服务模块。

package categoryv1

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
	CategoryService_CreateCategory_FullMethodName         = "/ecommerce.category.v1.CategoryService/CreateCategory"
	CategoryService_GetLeafCategories_FullMethodName      = "/ecommerce.category.v1.CategoryService/GetLeafCategories"
	CategoryService_BatchGetCategories_FullMethodName     = "/ecommerce.category.v1.CategoryService/BatchGetCategories"
	CategoryService_GetCategory_FullMethodName            = "/ecommerce.category.v1.CategoryService/GetCategory"
	CategoryService_UpdateCategory_FullMethodName         = "/ecommerce.category.v1.CategoryService/UpdateCategory"
	CategoryService_DeleteCategory_FullMethodName         = "/ecommerce.category.v1.CategoryService/DeleteCategory"
	CategoryService_GetSubTree_FullMethodName             = "/ecommerce.category.v1.CategoryService/GetSubTree"
	CategoryService_GetDirectSubCategories_FullMethodName = "/ecommerce.category.v1.CategoryService/GetDirectSubCategories"
	CategoryService_GetCategoryPath_FullMethodName        = "/ecommerce.category.v1.CategoryService/GetCategoryPath"
	CategoryService_GetClosureRelations_FullMethodName    = "/ecommerce.category.v1.CategoryService/GetClosureRelations"
	CategoryService_UpdateClosureDepth_FullMethodName     = "/ecommerce.category.v1.CategoryService/UpdateClosureDepth"
)

// CategoryServiceClient is the client API for CategoryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// 分类服务
type CategoryServiceClient interface {
	// 创建分类
	CreateCategory(ctx context.Context, in *CreateCategoryRequest, opts ...grpc.CallOption) (*Category, error)
	// 获取所有叶子节点
	GetLeafCategories(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Categories, error)
	// 批量查询分类
	BatchGetCategories(ctx context.Context, in *BatchGetCategoriesRequest, opts ...grpc.CallOption) (*Categories, error)
	// 获取单个分类
	GetCategory(ctx context.Context, in *GetCategoryRequest, opts ...grpc.CallOption) (*Category, error)
	// 更新分类
	UpdateCategory(ctx context.Context, in *UpdateCategoryRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// 删除分类
	DeleteCategory(ctx context.Context, in *DeleteCategoryRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// 获取子树
	GetSubTree(ctx context.Context, in *GetSubTreeRequest, opts ...grpc.CallOption) (*Categories, error)
	// 获取直接子分类（只返回下一级）
	GetDirectSubCategories(ctx context.Context, in *GetDirectSubCategoriesRequest, opts ...grpc.CallOption) (*Categories, error)
	// 获取分类路径
	GetCategoryPath(ctx context.Context, in *GetCategoryPathRequest, opts ...grpc.CallOption) (*Categories, error)
	// 获取闭包关系
	GetClosureRelations(ctx context.Context, in *GetClosureRequest, opts ...grpc.CallOption) (*ClosureRelations, error)
	// 更新闭包关系深度
	UpdateClosureDepth(ctx context.Context, in *UpdateClosureDepthRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type categoryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCategoryServiceClient(cc grpc.ClientConnInterface) CategoryServiceClient {
	return &categoryServiceClient{cc}
}

func (c *categoryServiceClient) CreateCategory(ctx context.Context, in *CreateCategoryRequest, opts ...grpc.CallOption) (*Category, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Category)
	err := c.cc.Invoke(ctx, CategoryService_CreateCategory_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *categoryServiceClient) GetLeafCategories(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Categories, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Categories)
	err := c.cc.Invoke(ctx, CategoryService_GetLeafCategories_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *categoryServiceClient) BatchGetCategories(ctx context.Context, in *BatchGetCategoriesRequest, opts ...grpc.CallOption) (*Categories, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Categories)
	err := c.cc.Invoke(ctx, CategoryService_BatchGetCategories_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *categoryServiceClient) GetCategory(ctx context.Context, in *GetCategoryRequest, opts ...grpc.CallOption) (*Category, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Category)
	err := c.cc.Invoke(ctx, CategoryService_GetCategory_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *categoryServiceClient) UpdateCategory(ctx context.Context, in *UpdateCategoryRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, CategoryService_UpdateCategory_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *categoryServiceClient) DeleteCategory(ctx context.Context, in *DeleteCategoryRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, CategoryService_DeleteCategory_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *categoryServiceClient) GetSubTree(ctx context.Context, in *GetSubTreeRequest, opts ...grpc.CallOption) (*Categories, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Categories)
	err := c.cc.Invoke(ctx, CategoryService_GetSubTree_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *categoryServiceClient) GetDirectSubCategories(ctx context.Context, in *GetDirectSubCategoriesRequest, opts ...grpc.CallOption) (*Categories, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Categories)
	err := c.cc.Invoke(ctx, CategoryService_GetDirectSubCategories_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *categoryServiceClient) GetCategoryPath(ctx context.Context, in *GetCategoryPathRequest, opts ...grpc.CallOption) (*Categories, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Categories)
	err := c.cc.Invoke(ctx, CategoryService_GetCategoryPath_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *categoryServiceClient) GetClosureRelations(ctx context.Context, in *GetClosureRequest, opts ...grpc.CallOption) (*ClosureRelations, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ClosureRelations)
	err := c.cc.Invoke(ctx, CategoryService_GetClosureRelations_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *categoryServiceClient) UpdateClosureDepth(ctx context.Context, in *UpdateClosureDepthRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, CategoryService_UpdateClosureDepth_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CategoryServiceServer is the server API for CategoryService service.
// All implementations must embed UnimplementedCategoryServiceServer
// for forward compatibility.
//
// 分类服务
type CategoryServiceServer interface {
	// 创建分类
	CreateCategory(context.Context, *CreateCategoryRequest) (*Category, error)
	// 获取所有叶子节点
	GetLeafCategories(context.Context, *emptypb.Empty) (*Categories, error)
	// 批量查询分类
	BatchGetCategories(context.Context, *BatchGetCategoriesRequest) (*Categories, error)
	// 获取单个分类
	GetCategory(context.Context, *GetCategoryRequest) (*Category, error)
	// 更新分类
	UpdateCategory(context.Context, *UpdateCategoryRequest) (*emptypb.Empty, error)
	// 删除分类
	DeleteCategory(context.Context, *DeleteCategoryRequest) (*emptypb.Empty, error)
	// 获取子树
	GetSubTree(context.Context, *GetSubTreeRequest) (*Categories, error)
	// 获取直接子分类（只返回下一级）
	GetDirectSubCategories(context.Context, *GetDirectSubCategoriesRequest) (*Categories, error)
	// 获取分类路径
	GetCategoryPath(context.Context, *GetCategoryPathRequest) (*Categories, error)
	// 获取闭包关系
	GetClosureRelations(context.Context, *GetClosureRequest) (*ClosureRelations, error)
	// 更新闭包关系深度
	UpdateClosureDepth(context.Context, *UpdateClosureDepthRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedCategoryServiceServer()
}

// UnimplementedCategoryServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCategoryServiceServer struct{}

func (UnimplementedCategoryServiceServer) CreateCategory(context.Context, *CreateCategoryRequest) (*Category, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCategory not implemented")
}
func (UnimplementedCategoryServiceServer) GetLeafCategories(context.Context, *emptypb.Empty) (*Categories, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLeafCategories not implemented")
}
func (UnimplementedCategoryServiceServer) BatchGetCategories(context.Context, *BatchGetCategoriesRequest) (*Categories, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchGetCategories not implemented")
}
func (UnimplementedCategoryServiceServer) GetCategory(context.Context, *GetCategoryRequest) (*Category, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCategory not implemented")
}
func (UnimplementedCategoryServiceServer) UpdateCategory(context.Context, *UpdateCategoryRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCategory not implemented")
}
func (UnimplementedCategoryServiceServer) DeleteCategory(context.Context, *DeleteCategoryRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCategory not implemented")
}
func (UnimplementedCategoryServiceServer) GetSubTree(context.Context, *GetSubTreeRequest) (*Categories, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSubTree not implemented")
}
func (UnimplementedCategoryServiceServer) GetDirectSubCategories(context.Context, *GetDirectSubCategoriesRequest) (*Categories, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDirectSubCategories not implemented")
}
func (UnimplementedCategoryServiceServer) GetCategoryPath(context.Context, *GetCategoryPathRequest) (*Categories, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCategoryPath not implemented")
}
func (UnimplementedCategoryServiceServer) GetClosureRelations(context.Context, *GetClosureRequest) (*ClosureRelations, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClosureRelations not implemented")
}
func (UnimplementedCategoryServiceServer) UpdateClosureDepth(context.Context, *UpdateClosureDepthRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateClosureDepth not implemented")
}
func (UnimplementedCategoryServiceServer) mustEmbedUnimplementedCategoryServiceServer() {}
func (UnimplementedCategoryServiceServer) testEmbeddedByValue()                         {}

// UnsafeCategoryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CategoryServiceServer will
// result in compilation errors.
type UnsafeCategoryServiceServer interface {
	mustEmbedUnimplementedCategoryServiceServer()
}

func RegisterCategoryServiceServer(s grpc.ServiceRegistrar, srv CategoryServiceServer) {
	// If the following call pancis, it indicates UnimplementedCategoryServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CategoryService_ServiceDesc, srv)
}

func _CategoryService_CreateCategory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCategoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CategoryServiceServer).CreateCategory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CategoryService_CreateCategory_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CategoryServiceServer).CreateCategory(ctx, req.(*CreateCategoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CategoryService_GetLeafCategories_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CategoryServiceServer).GetLeafCategories(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CategoryService_GetLeafCategories_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CategoryServiceServer).GetLeafCategories(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _CategoryService_BatchGetCategories_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchGetCategoriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CategoryServiceServer).BatchGetCategories(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CategoryService_BatchGetCategories_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CategoryServiceServer).BatchGetCategories(ctx, req.(*BatchGetCategoriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CategoryService_GetCategory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCategoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CategoryServiceServer).GetCategory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CategoryService_GetCategory_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CategoryServiceServer).GetCategory(ctx, req.(*GetCategoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CategoryService_UpdateCategory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCategoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CategoryServiceServer).UpdateCategory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CategoryService_UpdateCategory_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CategoryServiceServer).UpdateCategory(ctx, req.(*UpdateCategoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CategoryService_DeleteCategory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCategoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CategoryServiceServer).DeleteCategory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CategoryService_DeleteCategory_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CategoryServiceServer).DeleteCategory(ctx, req.(*DeleteCategoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CategoryService_GetSubTree_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSubTreeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CategoryServiceServer).GetSubTree(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CategoryService_GetSubTree_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CategoryServiceServer).GetSubTree(ctx, req.(*GetSubTreeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CategoryService_GetDirectSubCategories_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDirectSubCategoriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CategoryServiceServer).GetDirectSubCategories(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CategoryService_GetDirectSubCategories_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CategoryServiceServer).GetDirectSubCategories(ctx, req.(*GetDirectSubCategoriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CategoryService_GetCategoryPath_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCategoryPathRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CategoryServiceServer).GetCategoryPath(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CategoryService_GetCategoryPath_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CategoryServiceServer).GetCategoryPath(ctx, req.(*GetCategoryPathRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CategoryService_GetClosureRelations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetClosureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CategoryServiceServer).GetClosureRelations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CategoryService_GetClosureRelations_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CategoryServiceServer).GetClosureRelations(ctx, req.(*GetClosureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CategoryService_UpdateClosureDepth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateClosureDepthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CategoryServiceServer).UpdateClosureDepth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CategoryService_UpdateClosureDepth_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CategoryServiceServer).UpdateClosureDepth(ctx, req.(*UpdateClosureDepthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CategoryService_ServiceDesc is the grpc.ServiceDesc for CategoryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CategoryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ecommerce.category.v1.CategoryService",
	HandlerType: (*CategoryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateCategory",
			Handler:    _CategoryService_CreateCategory_Handler,
		},
		{
			MethodName: "GetLeafCategories",
			Handler:    _CategoryService_GetLeafCategories_Handler,
		},
		{
			MethodName: "BatchGetCategories",
			Handler:    _CategoryService_BatchGetCategories_Handler,
		},
		{
			MethodName: "GetCategory",
			Handler:    _CategoryService_GetCategory_Handler,
		},
		{
			MethodName: "UpdateCategory",
			Handler:    _CategoryService_UpdateCategory_Handler,
		},
		{
			MethodName: "DeleteCategory",
			Handler:    _CategoryService_DeleteCategory_Handler,
		},
		{
			MethodName: "GetSubTree",
			Handler:    _CategoryService_GetSubTree_Handler,
		},
		{
			MethodName: "GetDirectSubCategories",
			Handler:    _CategoryService_GetDirectSubCategories_Handler,
		},
		{
			MethodName: "GetCategoryPath",
			Handler:    _CategoryService_GetCategoryPath_Handler,
		},
		{
			MethodName: "GetClosureRelations",
			Handler:    _CategoryService_GetClosureRelations_Handler,
		},
		{
			MethodName: "UpdateClosureDepth",
			Handler:    _CategoryService_UpdateClosureDepth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/category.proto",
}
