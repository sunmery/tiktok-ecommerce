// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.3
// - protoc             v5.29.3
// source: v1/product.proto

package productv1

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

const OperationProductServiceAuditProduct = "/ecommerce.product.v1.ProductService/AuditProduct"
const OperationProductServiceCreateProduct = "/ecommerce.product.v1.ProductService/CreateProduct"
const OperationProductServiceDeleteProduct = "/ecommerce.product.v1.ProductService/DeleteProduct"
const OperationProductServiceGetMerchantProducts = "/ecommerce.product.v1.ProductService/GetMerchantProducts"
const OperationProductServiceGetProduct = "/ecommerce.product.v1.ProductService/GetProduct"
const OperationProductServiceListProductsByCategory = "/ecommerce.product.v1.ProductService/ListProductsByCategory"
const OperationProductServiceListRandomProducts = "/ecommerce.product.v1.ProductService/ListRandomProducts"
const OperationProductServiceSearchProductsByName = "/ecommerce.product.v1.ProductService/SearchProductsByName"
const OperationProductServiceSubmitForAudit = "/ecommerce.product.v1.ProductService/SubmitForAudit"
const OperationProductServiceUpdateProduct = "/ecommerce.product.v1.ProductService/UpdateProduct"

type ProductServiceHTTPServer interface {
	// AuditProduct 审核商品
	AuditProduct(context.Context, *AuditProductRequest) (*AuditRecord, error)
	// CreateProduct 创建商品（草稿状态）
	CreateProduct(context.Context, *CreateProductRequest) (*CreateProductReply, error)
	// DeleteProduct 删除商品（软删除）
	DeleteProduct(context.Context, *DeleteProductRequest) (*emptypb.Empty, error)
	// GetMerchantProducts 获取商家对应的商品
	GetMerchantProducts(context.Context, *GetMerchantProductRequest) (*Products, error)
	// GetProduct 获取商品详情
	GetProduct(context.Context, *GetProductRequest) (*Product, error)
	// ListProductsByCategory 根据商品分类查询
	ListProductsByCategory(context.Context, *ListProductsByCategoryRequest) (*Products, error)
	// ListRandomProducts 随机返回商品数据
	ListRandomProducts(context.Context, *ListRandomProductsRequest) (*Products, error)
	// SearchProductsByName 根据商品名称模糊查询
	SearchProductsByName(context.Context, *SearchProductRequest) (*Products, error)
	// SubmitForAudit 提交商品审核
	SubmitForAudit(context.Context, *SubmitAuditRequest) (*AuditRecord, error)
	// UpdateProduct 更新商品信息
	UpdateProduct(context.Context, *UpdateProductRequest) (*Product, error)
}

func RegisterProductServiceHTTPServer(s *http.Server, srv ProductServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/v1/products", _ProductService_CreateProduct0_HTTP_Handler(srv))
	r.PUT("/v1/products/{id}", _ProductService_UpdateProduct0_HTTP_Handler(srv))
	r.POST("/v1/products/{product_id}/submit-audit", _ProductService_SubmitForAudit0_HTTP_Handler(srv))
	r.POST("/v1/products/{product_id}/audit", _ProductService_AuditProduct0_HTTP_Handler(srv))
	r.GET("/v1/products", _ProductService_ListRandomProducts0_HTTP_Handler(srv))
	r.GET("/v1/products/{id}", _ProductService_GetProduct0_HTTP_Handler(srv))
	r.GET("/v1/merchants/products", _ProductService_GetMerchantProducts0_HTTP_Handler(srv))
	r.GET("/v1/products/{name}", _ProductService_SearchProductsByName0_HTTP_Handler(srv))
	r.GET("/v1/products/categories/{name}", _ProductService_ListProductsByCategory0_HTTP_Handler(srv))
	r.DELETE("/v1/products", _ProductService_DeleteProduct0_HTTP_Handler(srv))
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
		reply := out.(*CreateProductReply)
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

func _ProductService_ListRandomProducts0_HTTP_Handler(srv ProductServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListRandomProductsRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationProductServiceListRandomProducts)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListRandomProducts(ctx, req.(*ListRandomProductsRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Products)
		return ctx.Result(200, reply)
	}
}

func _ProductService_GetProduct0_HTTP_Handler(srv ProductServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetProductRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationProductServiceGetProduct)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetProduct(ctx, req.(*GetProductRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Product)
		return ctx.Result(200, reply)
	}
}

func _ProductService_GetMerchantProducts0_HTTP_Handler(srv ProductServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetMerchantProductRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationProductServiceGetMerchantProducts)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetMerchantProducts(ctx, req.(*GetMerchantProductRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Products)
		return ctx.Result(200, reply)
	}
}

func _ProductService_SearchProductsByName0_HTTP_Handler(srv ProductServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in SearchProductRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationProductServiceSearchProductsByName)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SearchProductsByName(ctx, req.(*SearchProductRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Products)
		return ctx.Result(200, reply)
	}
}

func _ProductService_ListProductsByCategory0_HTTP_Handler(srv ProductServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListProductsByCategoryRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationProductServiceListProductsByCategory)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListProductsByCategory(ctx, req.(*ListProductsByCategoryRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*Products)
		return ctx.Result(200, reply)
	}
}

func _ProductService_DeleteProduct0_HTTP_Handler(srv ProductServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteProductRequest
		if err := ctx.BindQuery(&in); err != nil {
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
	CreateProduct(ctx context.Context, req *CreateProductRequest, opts ...http.CallOption) (rsp *CreateProductReply, err error)
	DeleteProduct(ctx context.Context, req *DeleteProductRequest, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	GetMerchantProducts(ctx context.Context, req *GetMerchantProductRequest, opts ...http.CallOption) (rsp *Products, err error)
	GetProduct(ctx context.Context, req *GetProductRequest, opts ...http.CallOption) (rsp *Product, err error)
	ListProductsByCategory(ctx context.Context, req *ListProductsByCategoryRequest, opts ...http.CallOption) (rsp *Products, err error)
	ListRandomProducts(ctx context.Context, req *ListRandomProductsRequest, opts ...http.CallOption) (rsp *Products, err error)
	SearchProductsByName(ctx context.Context, req *SearchProductRequest, opts ...http.CallOption) (rsp *Products, err error)
	SubmitForAudit(ctx context.Context, req *SubmitAuditRequest, opts ...http.CallOption) (rsp *AuditRecord, err error)
	UpdateProduct(ctx context.Context, req *UpdateProductRequest, opts ...http.CallOption) (rsp *Product, err error)
}

type ProductServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewProductServiceHTTPClient(client *http.Client) ProductServiceHTTPClient {
	return &ProductServiceHTTPClientImpl{client}
}

func (c *ProductServiceHTTPClientImpl) AuditProduct(ctx context.Context, in *AuditProductRequest, opts ...http.CallOption) (*AuditRecord, error) {
	var out AuditRecord
	pattern := "/v1/products/{product_id}/audit"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationProductServiceAuditProduct))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ProductServiceHTTPClientImpl) CreateProduct(ctx context.Context, in *CreateProductRequest, opts ...http.CallOption) (*CreateProductReply, error) {
	var out CreateProductReply
	pattern := "/v1/products"
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
	pattern := "/v1/products"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationProductServiceDeleteProduct))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ProductServiceHTTPClientImpl) GetMerchantProducts(ctx context.Context, in *GetMerchantProductRequest, opts ...http.CallOption) (*Products, error) {
	var out Products
	pattern := "/v1/merchants/products"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationProductServiceGetMerchantProducts))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
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

func (c *ProductServiceHTTPClientImpl) ListProductsByCategory(ctx context.Context, in *ListProductsByCategoryRequest, opts ...http.CallOption) (*Products, error) {
	var out Products
	pattern := "/v1/products/categories/{name}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationProductServiceListProductsByCategory))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ProductServiceHTTPClientImpl) ListRandomProducts(ctx context.Context, in *ListRandomProductsRequest, opts ...http.CallOption) (*Products, error) {
	var out Products
	pattern := "/v1/products"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationProductServiceListRandomProducts))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *ProductServiceHTTPClientImpl) SearchProductsByName(ctx context.Context, in *SearchProductRequest, opts ...http.CallOption) (*Products, error) {
	var out Products
	pattern := "/v1/products/{name}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationProductServiceSearchProductsByName))
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
