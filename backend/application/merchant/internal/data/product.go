package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"

	"backend/application/merchant/internal/data/models"

	"backend/constants"

	"github.com/go-kratos/kratos/v2/metadata"

	productv1 "backend/api/product/v1"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/jackc/pgx/v5"

	"backend/pkg/types"

	v1 "backend/api/product/v1"
	"backend/application/merchant/internal/biz"
)

type productRepo struct {
	data *Data
	log  *log.Helper
}

func NewProductRepo(data *Data, logger log.Logger) biz.ProductRepo {
	return &productRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (p *productRepo) GetMerchantProducts(ctx context.Context, req *biz.GetMerchantProducts) (*biz.Products, error) {
	db := p.data.DB(ctx)

	// 获取基础信息
	merchantID := types.ToPgUUID(req.MerchantID)
	merchantProducts, err := db.GetMerchantProducts(ctx, models.GetMerchantProductsParams{
		MerchantID: merchantID,
		Page:       &req.Page,
		Pagesize:   &req.PageSize,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, v1.ErrorInvalidStatus("GetMerchantProducts 内部错误")
	}

	// 收集所有不同的分类ID
	categoryIDs := make([]int64, 0)
	categoryIDMap := make(map[int64]bool)
	for _, product := range merchantProducts {
		if !categoryIDMap[product.CategoryID] {
			categoryIDMap[product.CategoryID] = true
			categoryIDs = append(categoryIDs, product.CategoryID)
		}
	}

	// 从分类服务获取分类信息
	// var categoryMap map[int64]*category.Category
	// if len(categoryIDs) > 0 {
	// 	categoriesResp, err := p.data.BatchGetCategories(ctx, &category.BatchGetCategoriesRequest{
	// 		Ids: categoryIDs,
	// 	})
	// 	if err != nil {
	// 		p.log.WithContext(ctx).Warnf("failed to get categories: %v", err)
	// 	} else {
	// 		categoryMap = make(map[int64]*category.Category)
	// 		for _, cat := range categoriesResp.Categories {
	// 			categoryMap[cat.Id] = cat
	// 		}
	// 	}
	// }

	var products []*biz.Product
	for _, product := range merchantProducts {
		// 处理价格
		price, err := types.NumericToFloat(product.Price.(pgtype.Numeric))
		if err != nil {
			p.log.WithContext(ctx).Warnf("unmarshal price error: %v", err)
		}

		// 处理图片
		var images []*biz.ProductImage
		if len(product.Images) > 0 {
			if err := json.Unmarshal(product.Images, &images); err != nil {
				p.log.WithContext(ctx).Warnf("unmarshal images error: %v", err)
				continue
			}
		}

		// 处理属性
		var attributes map[string]*biz.AttributeValue
		if len(product.Attributes) > 0 {
			var rawJSON map[string]any
			if err := json.Unmarshal(product.Attributes, &rawJSON); err != nil {
				p.log.WithContext(ctx).Warnf("unmarshal attributes error: %v", err)
				continue
			} else {
				attributes = make(map[string]*biz.AttributeValue)
				for key, value := range rawJSON {
					switch v := value.(type) {
					case string:
						attributes[key] = &biz.AttributeValue{StringValue: v}
					case []any:
						items := make([]string, len(v))
						for i, item := range v {
							items[i] = item.(string)
						}
						attributes[key] = &biz.AttributeValue{
							ArrayValue: &biz.ArrayValue{Items: items},
						}
					case map[string]any:
						fields := make(map[string]*biz.AttributeValue)
						for k, val := range v {
							fields[k] = &biz.AttributeValue{StringValue: val.(string)}
						}
						attributes[key] = &biz.AttributeValue{
							ObjectValue: &biz.NestedObject{Fields: fields},
						}
					}
				}
			}
		}

		// 构建分类信息
		categoryInfo := biz.CategoryInfo{
			CategoryId: uint64(product.CategoryID),
		}

		// 如果找到了分类信息，则设置分类名称
		// if c, ok := categoryMap[product.CategoryID]; ok {
		// 	categoryInfo.CategoryName = c.Name
		// 	categoryInfo.SortOrder = c.SortOrder
		// }

		products = append(products, &biz.Product{
			ID:          product.ID,
			MerchantId:  product.MerchantID,
			Name:        product.Name,
			Price:       price,
			Description: *product.Description,
			Images:      images,
			Status:      constants.ProductStatus(product.Status),
			Category:    categoryInfo,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
			Attributes:  attributes,
			Inventory: biz.Inventory{
				ProductId:  product.ID,
				MerchantId: product.MerchantID,
				Stock:      product.Stock,
			},
		})
	}

	return &biz.Products{Items: products}, nil
}

func (p *productRepo) UpdateProduct(ctx context.Context, req *biz.UpdateProductRequest) (*biz.UpdateProductReply, error) {
	db := p.data.DB(ctx)

	ctx = metadata.AppendToClientContext(ctx, constants.UserId, req.MerchantID.String())
	oldProduct, err := p.data.productv1.GetProduct(ctx, &productv1.GetProductRequest{
		Id:         req.ID.String(),
		MerchantId: req.MerchantID.String(),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, v1.ErrorProductNotFound("查询不到该商品")
		}
		return nil, v1.ErrorInvalidStatus("failed to get product: %w", err)
	}

	// 准备更新参数
	productId := types.ToPgUUID(req.ID)
	merchantId := types.ToPgUUID(req.MerchantID)
	price, err := types.Float64ToNumeric(*req.Price)
	if err != nil {
		return nil, fmt.Errorf("invalid price: %w", err)
	}
	attributes, err := json.Marshal(req.Attributes)
	if err != nil {
		return nil, fmt.Errorf("invalid attributes: %w", err)
	}
	status := int16(req.Status)
	params := models.UpdateProductParams{
		Stock:       &req.Stock,
		Name:        req.Name,
		Description: req.Description,
		Price:       price,
		Status:      &status,
		ProductID:   productId,
		MerchantID:  merchantId,
		Attributes:  attributes,
		Url:         &req.Url,
	}

	// 执行更新
	err = db.UpdateProduct(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return &biz.UpdateProductReply{
		Code:    200,
		Message: fmt.Sprintf("更新商品'%s'成功", oldProduct.Name),
	}, nil
}
