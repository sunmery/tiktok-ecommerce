package service

import (
	"context"

	productv1 "backend/api/product/v1"

	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "backend/api/merchant/product/v1"
	"backend/application/merchant/internal/biz"
	"backend/pkg"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (uc *ProductService) GetMerchantProducts(ctx context.Context, _ *v1.GetMerchantProductRequest) (*v1.Products, error) {
	merchantId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid merchantId ID")
	}
	products, err := uc.pc.GetMerchantProducts(ctx, &biz.GetMerchantProducts{
		MerchantID: merchantId,
	})
	if err != nil {
		return nil, err
	}

	var pbProducts []*productv1.Product
	for _, product := range products.Items {
		pbProducts = append(pbProducts, convertBizProductToPB(product))
	}
	return &v1.Products{
		Items: pbProducts,
	}, nil
}

// 辅助转换方法
func convertBizProductToPB(p *biz.Product) *productv1.Product {
	pbProduct := &productv1.Product{
		Id:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Status:      convertBizStatusToPB(p.Status),
		MerchantId:  p.MerchantId.String(),
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
		Category: &productv1.CategoryInfo{
			CategoryId:   uint32(p.Category.CategoryId),
			CategoryName: p.Category.CategoryName,
		},
	}

	for _, img := range p.Images {
		// 安全转换 SortOrder
		var sortOrder int32
		if img.SortOrder != nil {
			sortOrder = int32(*img.SortOrder) // 解引用并转换类型
		} else {
			sortOrder = 0 // 默认值
		}
		pbProduct.Images = append(pbProduct.Images, &productv1.Image{
			Url:       img.URL,
			IsPrimary: img.IsPrimary,
			SortOrder: sortOrder,
		})
	}

	return pbProduct
}

func convertBizStatusToPB(s biz.ProductStatus) productv1.ProductStatus {
	switch s {
	case biz.ProductStatusDraft:
		return productv1.ProductStatus_PRODUCT_STATUS_DRAFT
	case biz.ProductStatusPending:
		return productv1.ProductStatus_PRODUCT_STATUS_PENDING
	case biz.ProductStatusApproved:
		return productv1.ProductStatus_PRODUCT_STATUS_APPROVED
	case biz.ProductStatusRejected:
		return productv1.ProductStatus_PRODUCT_STATUS_REJECTED
	default:
		return productv1.ProductStatus_PRODUCT_STATUS_DRAFT
	}
}
