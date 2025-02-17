// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.3
// - protoc             v5.29.3
// source: api/product/v1/service.proto

package product

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

<<<<<<< HEAD
const OperationProductCatalogServiceCreateCategory = "/product.service.v1.ProductCatalogService/CreateCategory"
const OperationProductCatalogServiceCreateProduct = "/product.service.v1.ProductCatalogService/CreateProduct"
const OperationProductCatalogServiceDeleteProduct = "/product.service.v1.ProductCatalogService/DeleteProduct"
const OperationProductCatalogServiceGetCategoryChildren = "/product.service.v1.ProductCatalogService/GetCategoryChildren"
const OperationProductCatalogServiceGetProduct = "/product.service.v1.ProductCatalogService/GetProduct"
const OperationProductCatalogServiceListCategories = "/product.service.v1.ProductCatalogService/ListCategories"
const OperationProductCatalogServiceListProducts = "/product.service.v1.ProductCatalogService/ListProducts"
const OperationProductCatalogServiceSearchProducts = "/product.service.v1.ProductCatalogService/SearchProducts"
const OperationProductCatalogServiceUpdateProduct = "/product.service.v1.ProductCatalogService/UpdateProduct"

type ProductCatalogServiceHTTPServer interface {
	CreateCategory(context.Context, *CreateCategoryReq) (*CategoryReply, error)
	CreateProduct(context.Context, *CreateProductRequest) (*ProductReply, error)
	DeleteProduct(context.Context, *DeleteProductReq) (*ProductReply, error)
	GetCategoryChildren(context.Context, *GetCategoryChildrenReq) (*GetCategoryChildrenResp, error)
	GetProduct(context.Context, *GetProductReq) (*ProductReply, error)
	ListCategories(context.Context, *emptypb.Empty) (*ListCategoriesResp, error)
	ListProducts(context.Context, *ListProductsReq) (*ListProductsResp, error)
	SearchProducts(context.Context, *SearchProductsReq) (*SearchProductsResp, error)
	UpdateProduct(context.Context, *UpdateProductRequest) (*ProductReply, error)
=======
const OperationProductServiceAuditProduct = "/api.product.v1.ProductService/AuditProduct"
const OperationProductServiceCreateProduct = "/api.product.v1.ProductService/CreateProduct"
const OperationProductServiceDeleteProduct = "/api.product.v1.ProductService/DeleteProduct"
const OperationProductServiceGetProduct = "/api.product.v1.ProductService/GetProduct"
const OperationProductServiceSubmitForAudit = "/api.product.v1.ProductService/SubmitForAudit"
const OperationProductServiceUpdateProduct = "/api.product.v1.ProductService/UpdateProduct"

type ProductServiceHTTPServer interface {
	// AuditProduct 审核商品
	AuditProduct(context.Context, *AuditProductRequest) (*AuditRecord, error)
	// CreateProduct 创建商品（草稿状态）
	CreateProduct(context.Context, *CreateProductRequest) (*Product, error)
	// DeleteProduct 删除商品（软删除）
	DeleteProduct(context.Context, *DeleteProductRequest) (*emptypb.Empty, error)
	// GetProduct 获取商品详情
	GetProduct(context.Context, *GetProductRequest) (*Product, error)
	// SubmitForAudit 提交商品审核
	SubmitForAudit(context.Context, *SubmitAuditRequest) (*AuditRecord, error)
	// UpdateProduct 更新商品信息
	UpdateProduct(context.Context, *UpdateProductRequest) (*Product, error)
>>>>>>> main
}

func RegisterProductServiceHTTPServer(s *http.Server, srv ProductServiceHTTPServer) {
	r := s.Route("/")
<<<<<<< HEAD
	r.POST("/v1/product", _ProductCatalogService_CreateProduct0_HTTP_Handler(srv))
	r.PATCH("/v1/product", _ProductCatalogService_UpdateProduct0_HTTP_Handler(srv))
	r.GET("/v1/product/list", _ProductCatalogService_ListProducts0_HTTP_Handler(srv))
	r.GET("/v1/product/{id}", _ProductCatalogService_GetProduct0_HTTP_Handler(srv))
	r.GET("/v1/product/search/{query}", _ProductCatalogService_SearchProducts0_HTTP_Handler(srv))
	r.DELETE("/v1/product", _ProductCatalogService_DeleteProduct0_HTTP_Handler(srv))
	r.GET("/v1/categories/tree", _ProductCatalogService_ListCategories0_HTTP_Handler(srv))
	r.POST("/v1/categories", _ProductCatalogService_CreateCategory0_HTTP_Handler(srv))
	r.GET("/v1/categories/children/{id}", _ProductCatalogService_GetCategoryChildren0_HTTP_Handler(srv))
=======
	r.POST("/v1/products", _ProductService_CreateProduct0_HTTP_Handler(srv))
	r.PUT("/v1/products/{id}", _ProductService_UpdateProduct0_HTTP_Handler(srv))
	r.POST("/v1/products/{product_id}/submit-audit", _ProductService_SubmitForAudit0_HTTP_Handler(srv))
	r.POST("/v1/products/{product_id}/audit", _ProductService_AuditProduct0_HTTP_Handler(srv))
	r.GET("/v1/products/{id}", _ProductService_GetProduct0_HTTP_Handler(srv))
	r.DELETE("/v1/products/{id}", _ProductService_DeleteProduct0_HTTP_Handler(srv))
>>>>>>> main
}

func _ProductService_CreateProduct0_HTTP_Handler(srv ProductServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateProductRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationProductServiceCreateProduct)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateProduct(ctx, req.(*CreateProductRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Product)
		return ctx.Result(200, reply)
	}
}

func _ProductService_UpdateProduct0_HTTP_Handler(srv ProductServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateProductRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationProductServiceUpdateProduct)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateProduct(ctx, req.(*UpdateProductRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Product)
		return ctx.Result(200, reply)
	}
}

func _ProductService_SubmitForAudit0_HTTP_Handler(srv ProductServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in SubmitAuditRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationProductServiceSubmitForAudit)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SubmitForAudit(ctx, req.(*SubmitAuditRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*AuditRecord)
		return ctx.Result(200, reply)
	}
}

func _ProductService_AuditProduct0_HTTP_Handler(srv ProductServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in AuditProductRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationProductServiceAuditProduct)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.AuditProduct(ctx, req.(*AuditProductRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*AuditRecord)
		return ctx.Result(200, reply)
	}
}

<<<<<<< HEAD
func _ProductCatalogService_DeleteProduct0_HTTP_Handler(srv ProductCatalogServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteProductReq
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationProductCatalogServiceDeleteProduct)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteProduct(ctx, req.(*DeleteProductReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ProductReply)
		return ctx.Result(200, reply)
	}
}

func _ProductCatalogService_ListCategories0_HTTP_Handler(srv ProductCatalogServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in emptypb.Empty
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationProductCatalogServiceListCategories)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListCategories(ctx, req.(*emptypb.Empty))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListCategoriesResp)
		return ctx.Result(200, reply)
	}
}

func _ProductCatalogService_CreateCategory0_HTTP_Handler(srv ProductCatalogServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateCategoryReq
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationProductCatalogServiceCreateCategory)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateCategory(ctx, req.(*CreateCategoryReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CategoryReply)
		return ctx.Result(200, reply)
	}
}

func _ProductCatalogService_GetCategoryChildren0_HTTP_Handler(srv ProductCatalogServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetCategoryChildrenReq
=======
func _ProductService_GetProduct0_HTTP_Handler(srv ProductServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetProductRequest
>>>>>>> main
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
<<<<<<< HEAD
		http.SetOperation(ctx, OperationProductCatalogServiceGetCategoryChildren)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetCategoryChildren(ctx, req.(*GetCategoryChildrenReq))
=======
		http.SetOperation(ctx, OperationProductServiceGetProduct)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetProduct(ctx, req.(*GetProductRequest))
>>>>>>> main
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
<<<<<<< HEAD
		reply := out.(*GetCategoryChildrenResp)
		return ctx.Result(200, reply)
	}
}

type ProductCatalogServiceHTTPClient interface {
	CreateCategory(ctx context.Context, req *CreateCategoryReq, opts ...http.CallOption) (rsp *CategoryReply, err error)
	CreateProduct(ctx context.Context, req *CreateProductRequest, opts ...http.CallOption) (rsp *ProductReply, err error)
	DeleteProduct(ctx context.Context, req *DeleteProductReq, opts ...http.CallOption) (rsp *ProductReply, err error)
	GetCategoryChildren(ctx context.Context, req *GetCategoryChildrenReq, opts ...http.CallOption) (rsp *GetCategoryChildrenResp, err error)
	GetProduct(ctx context.Context, req *GetProductReq, opts ...http.CallOption) (rsp *ProductReply, err error)
	ListCategories(ctx context.Context, req *emptypb.Empty, opts ...http.CallOption) (rsp *ListCategoriesResp, err error)
	ListProducts(ctx context.Context, req *ListProductsReq, opts ...http.CallOption) (rsp *ListProductsResp, err error)
	SearchProducts(ctx context.Context, req *SearchProductsReq, opts ...http.CallOption) (rsp *SearchProductsResp, err error)
	UpdateProduct(ctx context.Context, req *UpdateProductRequest, opts ...http.CallOption) (rsp *ProductReply, err error)
=======
		reply := out.(*Product)
		return ctx.Result(200, reply)
	}
>>>>>>> main
}

func _ProductService_DeleteProduct0_HTTP_Handler(srv ProductServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteProductRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationProductServiceDeleteProduct)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteProduct(ctx, req.(*DeleteProductRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

type ProductServiceHTTPClient interface {
	AuditProduct(ctx context.Context, req *AuditProductRequest, opts ...http.CallOption) (rsp *AuditRecord, err error)
	CreateProduct(ctx context.Context, req *CreateProductRequest, opts ...http.CallOption) (rsp *Product, err error)
	DeleteProduct(ctx context.Context, req *DeleteProductRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	GetProduct(ctx context.Context, req *GetProductRequest, opts ...http.CallOption) (rsp *Product, err error)
	SubmitForAudit(ctx context.Context, req *SubmitAuditRequest, opts ...http.CallOption) (rsp *AuditRecord, err error)
	UpdateProduct(ctx context.Context, req *UpdateProductRequest, opts ...http.CallOption) (rsp *Product, err error)
}

type ProductServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewProductServiceHTTPClient(client *http.Client) ProductServiceHTTPClient {
	return &ProductServiceHTTPClientImpl{client}
}

<<<<<<< HEAD
func (c *ProductCatalogServiceHTTPClientImpl) CreateCategory(ctx context.Context, in *CreateCategoryReq, opts ...http.CallOption) (*CategoryReply, error) {
	var out CategoryReply
	pattern := "/v1/categories"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationProductCatalogServiceCreateCategory))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ProductCatalogServiceHTTPClientImpl) CreateProduct(ctx context.Context, in *CreateProductRequest, opts ...http.CallOption) (*ProductReply, error) {
	var out ProductReply
	pattern := "/v1/product"
=======
func (c *ProductServiceHTTPClientImpl) AuditProduct(ctx context.Context, in *AuditProductRequest, opts ...http.CallOption) (*AuditRecord, error) {
	var out AuditRecord
	pattern := "/v1/products/{product_id}/audit"
>>>>>>> main
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationProductServiceAuditProduct))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

<<<<<<< HEAD
func (c *ProductCatalogServiceHTTPClientImpl) DeleteProduct(ctx context.Context, in *DeleteProductReq, opts ...http.CallOption) (*ProductReply, error) {
	var out ProductReply
	pattern := "/v1/product"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationProductCatalogServiceDeleteProduct))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ProductCatalogServiceHTTPClientImpl) GetCategoryChildren(ctx context.Context, in *GetCategoryChildrenReq, opts ...http.CallOption) (*GetCategoryChildrenResp, error) {
	var out GetCategoryChildrenResp
	pattern := "/v1/categories/children/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationProductCatalogServiceGetCategoryChildren))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ProductCatalogServiceHTTPClientImpl) GetProduct(ctx context.Context, in *GetProductReq, opts ...http.CallOption) (*ProductReply, error) {
	var out ProductReply
	pattern := "/v1/product/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationProductCatalogServiceGetProduct))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ProductCatalogServiceHTTPClientImpl) ListCategories(ctx context.Context, in *emptypb.Empty, opts ...http.CallOption) (*ListCategoriesResp, error) {
	var out ListCategoriesResp
	pattern := "/v1/categories/tree"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationProductCatalogServiceListCategories))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ProductCatalogServiceHTTPClientImpl) ListProducts(ctx context.Context, in *ListProductsReq, opts ...http.CallOption) (*ListProductsResp, error) {
	var out ListProductsResp
	pattern := "/v1/product/list"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationProductCatalogServiceListProducts))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ProductCatalogServiceHTTPClientImpl) SearchProducts(ctx context.Context, in *SearchProductsReq, opts ...http.CallOption) (*SearchProductsResp, error) {
	var out SearchProductsResp
	pattern := "/v1/product/search/{query}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationProductCatalogServiceSearchProducts))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ProductCatalogServiceHTTPClientImpl) UpdateProduct(ctx context.Context, in *UpdateProductRequest, opts ...http.CallOption) (*ProductReply, error) {
	var out ProductReply
	pattern := "/v1/product"
=======
func (c *ProductServiceHTTPClientImpl) CreateProduct(ctx context.Context, in *CreateProductRequest, opts ...http.CallOption) (*Product, error) {
	var out Product
	pattern := "/v1/products"
>>>>>>> main
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationProductServiceCreateProduct))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ProductServiceHTTPClientImpl) DeleteProduct(ctx context.Context, in *DeleteProductRequest, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/v1/products/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationProductServiceDeleteProduct))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ProductServiceHTTPClientImpl) GetProduct(ctx context.Context, in *GetProductRequest, opts ...http.CallOption) (*Product, error) {
	var out Product
	pattern := "/v1/products/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationProductServiceGetProduct))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ProductServiceHTTPClientImpl) SubmitForAudit(ctx context.Context, in *SubmitAuditRequest, opts ...http.CallOption) (*AuditRecord, error) {
	var out AuditRecord
	pattern := "/v1/products/{product_id}/submit-audit"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationProductServiceSubmitForAudit))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ProductServiceHTTPClientImpl) UpdateProduct(ctx context.Context, in *UpdateProductRequest, opts ...http.CallOption) (*Product, error) {
	var out Product
	pattern := "/v1/products/{id}"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationProductServiceUpdateProduct))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
