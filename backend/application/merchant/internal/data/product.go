package data

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"backend/pkg/types"

	v1 "backend/api/product/v1"
	"backend/application/merchant/internal/biz"
)

func (p *productRepo) GetMerchantProducts(ctx context.Context, req *biz.GetMerchantProducts) (*biz.Products, error) {
	db := p.data.DB(ctx)

	// 获取基础信息
	merchantID := types.ToPgUUID(req.MerchantID)
	merchantProducts, err := db.GetMerchantProducts(ctx, merchantID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, v1.ErrorProductNotFound("查询不到商家的商品")
		}
		return nil, v1.ErrorInvalidStatus("GetMerchantProducts 内部错误")
	}
	var products []*biz.Product
	for _, product := range merchantProducts {
		products = append(products, &biz.Product{
			ID:         product.ID,
			MerchantId: product.MerchantID,
			Name:       product.Name,
			// Price:       price,
			Description: *product.Description,
			// Images:      convertImages(product.Images),
			Status:    biz.ProductStatus(product.Status),
			CreatedAt: product.CreatedAt,
			UpdatedAt: product.UpdatedAt,
			// Attributes:  product.Attributes,
		})
	}

	return &biz.Products{Items: products}, nil
}
