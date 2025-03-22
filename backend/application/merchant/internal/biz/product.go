package biz

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// GetMerchantProducts 获取商家自身商品列表请求
type GetMerchantProducts struct {
	MerchantID uuid.UUID
}

// 商品实体
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
		SortOrder    int32
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

func (uc *ProductUsecase) GetMerchantProducts(ctx context.Context, req *GetMerchantProducts) (*Products, error) {
	return uc.repo.GetMerchantProducts(ctx, req)
}
