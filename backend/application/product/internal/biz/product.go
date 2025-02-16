package biz

import (
	"context"
	"github.com/google/uuid"
	"time"
)

// Product 商品领域模型
type Product struct {
	ID          uuid.UUID
	Name        string
	Picture     string
	Price       float64
	Description string
	Stock       int32
	CategoryID  uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UpdateProductRequest 更新商品请求结构体
type UpdateProductRequest struct {
	Name        *string
	Price       *float64
	Picture     *string
	Description string
	Stock       *int
	Category    *string
}
type AddProductRequest struct {
}
type ListProductsReq struct {
	Page         uint   `json:"page"`
	PageSize     uint   `json:"pageSize"`
	CategoryName string `json:"categoryName"`
}

type ListProductsResp struct {
	Product []*Product `json:"product"`
}

type GetProductResp struct {
	Product *Product `json:"product"`
}

type SearchProductsReq struct {
	Query string `json:"query"`
}
type SearchProductsResp struct {
	Result []*Product `json:"result"`
}

// ProductRepo is a Greater repo.
type ProductRepo interface {
	CreateProduct(ctx context.Context, req *Product) (*Product, error)
	AddProduct(ctx context.Context, req *UpdateProductRequest) (*Product, error)
	ListProducts(ctx context.Context, req *ListProductsReq) (*ListProductsResp, error)
	GetProduct(ctx context.Context, id uint32) (*GetProductResp, error)
	SearchProducts(ctx context.Context, req *SearchProductsReq) (*SearchProductsResp, error)
}

// CreateProduct 创建商品
func (uc *ProductUsecase) CreateProduct(ctx context.Context, req *Product) (*Product, error) {
	uc.log.WithContext(ctx).Infof("CreateProduct: %v", req)
	return uc.repo.CreateProduct(ctx, req)
}

func (uc *ProductUsecase) ListProducts(ctx context.Context, req *ListProductsReq) (*ListProductsResp, error) {
	uc.log.WithContext(ctx).Infof("ListProducts: %v", req)
	return uc.repo.ListProducts(ctx, req)
}

func (uc *ProductUsecase) GetProduct(ctx context.Context, id uint32) (*GetProductResp, error) {
	uc.log.WithContext(ctx).Infof("GetProductReq: %v", id)
	return uc.repo.GetProduct(ctx, id)
}

func (uc *ProductUsecase) SearchProducts(ctx context.Context, req *SearchProductsReq) (*SearchProductsResp, error) {
	uc.log.WithContext(ctx).Infof("SearchProducts: %v", req)
	return uc.repo.SearchProducts(ctx, req)
}
