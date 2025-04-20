package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	productv1 "backend/api/product/v1"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	v1 "backend/api/merchant/product/v1"
	"backend/application/merchant/internal/biz"
	"backend/pkg"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductService struct {
	v1.UnimplementedProductServer
	pc *biz.ProductUsecase
}

func NewProductService(pc *biz.ProductUsecase) *ProductService {
	return &ProductService{pc: pc}
}

func (uc *ProductService) GetMerchantProducts(ctx context.Context, req *v1.GetMerchantProductRequest) (*productv1.Products, error) {
	merchantId, err := pkg.GetMetadataUesrID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid merchantId ID")
	}
	page := (req.Page - 1) * req.PageSize
	products, err := uc.pc.GetMerchantProducts(ctx, &biz.GetMerchantProducts{
		MerchantID: merchantId,
		Page:       int64(page),
		PageSize:   int64(req.PageSize),
	})
	if err != nil {
		return nil, err
	}

	var pbProducts []*productv1.Product
	for _, product := range products.Items {
		pbProducts = append(pbProducts, convertBizProductToPB(product))
	}
	return &productv1.Products{
		Items: pbProducts,
	}, nil
}

func (uc *ProductService) UpdateProduct(ctx context.Context, req *v1.UpdateProductRequest) (*v1.UpdateProductReply, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid productId")
	}
	merchantId, err := uuid.Parse(req.MerchantId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid merchantId")
	}

	updateReq := biz.UpdateProductRequest{
		Stock:       req.Stock,
		Url:         req.Url,
		Attributes:  parseProtoValue(req.Attributes),
		ID:          id,
		MerchantID:  merchantId,
		Name:        &req.Name,
		Price:       &req.Price,
		Description: &req.Description,
	}

	result, err := uc.pc.UpdateProduct(ctx, &updateReq)
	if err != nil {
		return nil, err
	}

	return &v1.UpdateProductReply{
		Message: result.Message,
		Code:    int32(result.Code),
	}, nil
}

// 辅助转换方法
func convertBizProductToPB(p *biz.Product) *productv1.Product {
	if p == nil {
		return nil
	}

	// 转换图片
	images := make([]*productv1.Image, 0)
	if p.Images != nil {
		for _, img := range p.Images {
			if img != nil {
				sortOrder := int32(0)
				if img.SortOrder != nil {
					sortOrder = int32(*img.SortOrder)
				}
				images = append(images, &productv1.Image{
					Url:       img.URL,
					IsPrimary: img.IsPrimary,
					SortOrder: sortOrder,
				})
			}
		}
	}

	// 转换商品属性
	var attributes *structpb.Value
	if p.Attributes != nil && len(p.Attributes) > 0 {
		// 将 biz.AttributeValue 转换为 map[string]interface{}
		rawMap := make(map[string]any)
		for key, value := range p.Attributes {
			if value.StringValue != "" {
				rawMap[key] = value.StringValue
			} else if value.ArrayValue != nil {
				rawMap[key] = value.ArrayValue.Items
			} else if value.ObjectValue != nil {
				objMap := make(map[string]any)
				for k, v := range value.ObjectValue.Fields {
					objMap[k] = v.StringValue
				}
				rawMap[key] = objMap
			}
		}

		// 转换为 structpb.Struct
		protoStruct, err := structpb.NewStruct(rawMap)
		if err != nil {
			log.Warn("Error creating struct: %v", err)
			attributes = nil
		} else {
			attributes = structpb.NewStructValue(protoStruct)
		}
	}

	// 构建返回结果
	result := &productv1.Product{
		Id:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Status:      uint32(p.Status),
		MerchantId:  p.MerchantId.String(),
		Images:      images,
		Attributes:  attributes,
		Category: &productv1.CategoryInfo{
			CategoryId:   uint32(p.Category.CategoryId),
			CategoryName: p.Category.CategoryName,
		},
		CreatedAt: timestamppb.New(p.CreatedAt),
		UpdatedAt: timestamppb.New(p.UpdatedAt),
	}

	// 添加库存信息
	if p.Inventory.ProductId != uuid.Nil {
		result.Inventory = &productv1.Inventory{
			ProductId:  p.Inventory.ProductId.String(),
			MerchantId: p.Inventory.MerchantId.String(),
			Stock:      uint32(p.Inventory.Stock),
		}
	}

	return result
}

func parseProtoValue(v *structpb.Value) map[string]any {
	if v == nil {
		return nil
	}
	if v.GetStructValue() == nil {
		return nil
	}
	return v.GetStructValue().AsMap()
}
