// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.4
// - protoc             v5.29.3
// source: v1/category.proto

package categoryv1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationCategoryServiceBatchGetCategories = "/ecommerce.category.v1.CategoryService/BatchGetCategories"
const OperationCategoryServiceCreateCategory = "/ecommerce.category.v1.CategoryService/CreateCategory"
const OperationCategoryServiceDeleteCategory = "/ecommerce.category.v1.CategoryService/DeleteCategory"
const OperationCategoryServiceGetCategory = "/ecommerce.category.v1.CategoryService/GetCategory"
const OperationCategoryServiceGetCategoryPath = "/ecommerce.category.v1.CategoryService/GetCategoryPath"
const OperationCategoryServiceGetClosureRelations = "/ecommerce.category.v1.CategoryService/GetClosureRelations"
const OperationCategoryServiceGetDirectSubCategories = "/ecommerce.category.v1.CategoryService/GetDirectSubCategories"
const OperationCategoryServiceGetLeafCategories = "/ecommerce.category.v1.CategoryService/GetLeafCategories"
const OperationCategoryServiceGetSubTree = "/ecommerce.category.v1.CategoryService/GetSubTree"
const OperationCategoryServiceUpdateCategory = "/ecommerce.category.v1.CategoryService/UpdateCategory"
const OperationCategoryServiceUpdateClosureDepth = "/ecommerce.category.v1.CategoryService/UpdateClosureDepth"

type CategoryServiceHTTPServer interface {
	// BatchGetCategories 批量查询分类
	BatchGetCategories(context.Context, *BatchGetCategoriesRequest) (*Categories, error)
	// CreateCategory 创建分类
	CreateCategory(context.Context, *CreateCategoryRequest) (*Category, error)
	// DeleteCategory 删除分类
	DeleteCategory(context.Context, *DeleteCategoryRequest) (*emptypb.Empty, error)
	// GetCategory 获取单个分类
	GetCategory(context.Context, *GetCategoryRequest) (*Category, error)
	// GetCategoryPath 获取分类路径
	GetCategoryPath(context.Context, *GetCategoryPathRequest) (*Categories, error)
	// GetClosureRelations 获取闭包关系
	GetClosureRelations(context.Context, *GetClosureRequest) (*ClosureRelations, error)
	// GetDirectSubCategories 获取直接子分类（只返回下一级）
	GetDirectSubCategories(context.Context, *GetDirectSubCategoriesRequest) (*Categories, error)
	// GetLeafCategories 获取所有叶子节点
	GetLeafCategories(context.Context, *emptypb.Empty) (*Categories, error)
	// GetSubTree 获取子树
	GetSubTree(context.Context, *GetSubTreeRequest) (*Categories, error)
	// UpdateCategory 更新分类
	UpdateCategory(context.Context, *UpdateCategoryRequest) (*emptypb.Empty, error)
	// UpdateClosureDepth 更新闭包关系深度
	UpdateClosureDepth(context.Context, *UpdateClosureDepthRequest) (*emptypb.Empty, error)
}

func RegisterCategoryServiceHTTPServer(s *http.Server, srv CategoryServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/v1/categories", _CategoryService_CreateCategory0_HTTP_Handler(srv))
	r.GET("/v1/categories/leaves", _CategoryService_GetLeafCategories0_HTTP_Handler(srv))
	r.GET("/v1/categories/batch", _CategoryService_BatchGetCategories0_HTTP_Handler(srv))
	r.GET("/v1/categories/{id}", _CategoryService_GetCategory0_HTTP_Handler(srv))
	r.PUT("/v1/categories/{id}", _CategoryService_UpdateCategory0_HTTP_Handler(srv))
	r.DELETE("/v1/categories/{id}", _CategoryService_DeleteCategory0_HTTP_Handler(srv))
	r.GET("/v1/categories/{root_id}/subtree", _CategoryService_GetSubTree0_HTTP_Handler(srv))
	r.GET("/v1/categories/{parent_id}/children", _CategoryService_GetDirectSubCategories0_HTTP_Handler(srv))
	r.GET("/v1/categories/{category_id}/path", _CategoryService_GetCategoryPath0_HTTP_Handler(srv))
	r.GET("/v1/categories/{category_id}/closure", _CategoryService_GetClosureRelations0_HTTP_Handler(srv))
	r.PATCH("/v1/categories/{category_id}/closure", _CategoryService_UpdateClosureDepth0_HTTP_Handler(srv))
}

func _CategoryService_CreateCategory0_HTTP_Handler(srv CategoryServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateCategoryRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationCategoryServiceCreateCategory)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateCategory(ctx, req.(*CreateCategoryRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Category)
		return ctx.Result(200, reply)
	}
}

func _CategoryService_GetLeafCategories0_HTTP_Handler(srv CategoryServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in emptypb.Empty
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationCategoryServiceGetLeafCategories)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetLeafCategories(ctx, req.(*emptypb.Empty))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Categories)
		return ctx.Result(200, reply)
	}
}

func _CategoryService_BatchGetCategories0_HTTP_Handler(srv CategoryServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in BatchGetCategoriesRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationCategoryServiceBatchGetCategories)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.BatchGetCategories(ctx, req.(*BatchGetCategoriesRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Categories)
		return ctx.Result(200, reply)
	}
}

func _CategoryService_GetCategory0_HTTP_Handler(srv CategoryServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetCategoryRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationCategoryServiceGetCategory)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetCategory(ctx, req.(*GetCategoryRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Category)
		return ctx.Result(200, reply)
	}
}

func _CategoryService_UpdateCategory0_HTTP_Handler(srv CategoryServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateCategoryRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationCategoryServiceUpdateCategory)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateCategory(ctx, req.(*UpdateCategoryRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _CategoryService_DeleteCategory0_HTTP_Handler(srv CategoryServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteCategoryRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationCategoryServiceDeleteCategory)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteCategory(ctx, req.(*DeleteCategoryRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _CategoryService_GetSubTree0_HTTP_Handler(srv CategoryServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetSubTreeRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationCategoryServiceGetSubTree)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetSubTree(ctx, req.(*GetSubTreeRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Categories)
		return ctx.Result(200, reply)
	}
}

func _CategoryService_GetDirectSubCategories0_HTTP_Handler(srv CategoryServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetDirectSubCategoriesRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationCategoryServiceGetDirectSubCategories)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetDirectSubCategories(ctx, req.(*GetDirectSubCategoriesRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Categories)
		return ctx.Result(200, reply)
	}
}

func _CategoryService_GetCategoryPath0_HTTP_Handler(srv CategoryServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetCategoryPathRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationCategoryServiceGetCategoryPath)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetCategoryPath(ctx, req.(*GetCategoryPathRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Categories)
		return ctx.Result(200, reply)
	}
}

func _CategoryService_GetClosureRelations0_HTTP_Handler(srv CategoryServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetClosureRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationCategoryServiceGetClosureRelations)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetClosureRelations(ctx, req.(*GetClosureRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ClosureRelations)
		return ctx.Result(200, reply)
	}
}

func _CategoryService_UpdateClosureDepth0_HTTP_Handler(srv CategoryServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateClosureDepthRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationCategoryServiceUpdateClosureDepth)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateClosureDepth(ctx, req.(*UpdateClosureDepthRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

type CategoryServiceHTTPClient interface {
	BatchGetCategories(ctx context.Context, req *BatchGetCategoriesRequest, opts ...http.CallOption) (rsp *Categories, err error)
	CreateCategory(ctx context.Context, req *CreateCategoryRequest, opts ...http.CallOption) (rsp *Category, err error)
	DeleteCategory(ctx context.Context, req *DeleteCategoryRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	GetCategory(ctx context.Context, req *GetCategoryRequest, opts ...http.CallOption) (rsp *Category, err error)
	GetCategoryPath(ctx context.Context, req *GetCategoryPathRequest, opts ...http.CallOption) (rsp *Categories, err error)
	GetClosureRelations(ctx context.Context, req *GetClosureRequest, opts ...http.CallOption) (rsp *ClosureRelations, err error)
	GetDirectSubCategories(ctx context.Context, req *GetDirectSubCategoriesRequest, opts ...http.CallOption) (rsp *Categories, err error)
	GetLeafCategories(ctx context.Context, req *emptypb.Empty, opts ...http.CallOption) (rsp *Categories, err error)
	GetSubTree(ctx context.Context, req *GetSubTreeRequest, opts ...http.CallOption) (rsp *Categories, err error)
	UpdateCategory(ctx context.Context, req *UpdateCategoryRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	UpdateClosureDepth(ctx context.Context, req *UpdateClosureDepthRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
}

type CategoryServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewCategoryServiceHTTPClient(client *http.Client) CategoryServiceHTTPClient {
	return &CategoryServiceHTTPClientImpl{client}
}

func (c *CategoryServiceHTTPClientImpl) BatchGetCategories(ctx context.Context, in *BatchGetCategoriesRequest, opts ...http.CallOption) (*Categories, error) {
	var out Categories
	pattern := "/v1/categories/batch"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationCategoryServiceBatchGetCategories))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *CategoryServiceHTTPClientImpl) CreateCategory(ctx context.Context, in *CreateCategoryRequest, opts ...http.CallOption) (*Category, error) {
	var out Category
	pattern := "/v1/categories"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationCategoryServiceCreateCategory))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *CategoryServiceHTTPClientImpl) DeleteCategory(ctx context.Context, in *DeleteCategoryRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/v1/categories/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationCategoryServiceDeleteCategory))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *CategoryServiceHTTPClientImpl) GetCategory(ctx context.Context, in *GetCategoryRequest, opts ...http.CallOption) (*Category, error) {
	var out Category
	pattern := "/v1/categories/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationCategoryServiceGetCategory))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *CategoryServiceHTTPClientImpl) GetCategoryPath(ctx context.Context, in *GetCategoryPathRequest, opts ...http.CallOption) (*Categories, error) {
	var out Categories
	pattern := "/v1/categories/{category_id}/path"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationCategoryServiceGetCategoryPath))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *CategoryServiceHTTPClientImpl) GetClosureRelations(ctx context.Context, in *GetClosureRequest, opts ...http.CallOption) (*ClosureRelations, error) {
	var out ClosureRelations
	pattern := "/v1/categories/{category_id}/closure"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationCategoryServiceGetClosureRelations))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *CategoryServiceHTTPClientImpl) GetDirectSubCategories(ctx context.Context, in *GetDirectSubCategoriesRequest, opts ...http.CallOption) (*Categories, error) {
	var out Categories
	pattern := "/v1/categories/{parent_id}/children"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationCategoryServiceGetDirectSubCategories))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *CategoryServiceHTTPClientImpl) GetLeafCategories(ctx context.Context, in *emptypb.Empty, opts ...http.CallOption) (*Categories, error) {
	var out Categories
	pattern := "/v1/categories/leaves"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationCategoryServiceGetLeafCategories))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *CategoryServiceHTTPClientImpl) GetSubTree(ctx context.Context, in *GetSubTreeRequest, opts ...http.CallOption) (*Categories, error) {
	var out Categories
	pattern := "/v1/categories/{root_id}/subtree"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationCategoryServiceGetSubTree))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *CategoryServiceHTTPClientImpl) UpdateCategory(ctx context.Context, in *UpdateCategoryRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/v1/categories/{id}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationCategoryServiceUpdateCategory))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *CategoryServiceHTTPClientImpl) UpdateClosureDepth(ctx context.Context, in *UpdateClosureDepthRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/v1/categories/{category_id}/closure"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationCategoryServiceUpdateClosureDepth))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PATCH", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
