package biz

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// GetMerchantProducts 获取商家自身商品列表请求
type GetMerchantProducts struct {
	MerchantID uuid.UUID
	Page       int64
	PageSize   int64
}

// Product 商品实体
type (
	ProductStatus uint

	ProductImage struct {
		URL       string
		IsPrimary bool
		SortOrder *int
	}
	CategoryInfo struct {
		CategoryId   uint64
		CategoryName string
	}
	AuditAction int

	ArrayValue struct {
		Items []string
	}

	NestedObject struct {
		Fields map[string]*AttributeValue
	}
	AttributeValue struct {
		StringValue string
		ArrayValue  *ArrayValue
		ObjectValue *NestedObject
	}
	// Product 商品领域模型
	Product struct {
		ID          uuid.UUID
		MerchantId  uuid.UUID
		Name        string
		Price       float64
		Description string
		Images      []*ProductImage
		Status      ProductStatus
		Category    CategoryInfo
		CreatedAt   time.Time
		UpdatedAt   time.Time
		Attributes  map[string]*AttributeValue
		Inventory   Inventory // 库存
	}
)

// Products 批量商品
type Products struct {
	Items []*Product
}

// UpdateProductRequest 更新商品请求结构体
type (
	UpdateProductRequest struct {
		ID          uuid.UUID
		MerchantID  uuid.UUID // 添加缺失字段
		Name        *string
		Price       *float64
		Description *string
	}
	UpdateProductReply struct {
		Code    uint
		Message string
	}
)

func (uc *ProductUsecase) GetMerchantProducts(ctx context.Context, req *GetMerchantProducts) (*Products, error) {
	return uc.repo.GetMerchantProducts(ctx, req)
}

func (uc *ProductUsecase) UpdateProduct(ctx context.Context, req *UpdateProductRequest) (*UpdateProductReply, error) {
	return uc.repo.UpdateProduct(ctx, req)
}
