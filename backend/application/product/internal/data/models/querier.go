// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package models

import (
	"context"
)

type Querier interface {
	// 批量插入图片
	//
	//  INSERT INTO products.product_images
	//      (merchant_id, product_id, url, is_primary, sort_order)
	//  SELECT m_id, p_id, u, is_p, s_ord
	//  FROM unnest(
	//               $1::bigint[],
	//               $2::bigint[],
	//               $3::text[],
	//               $4::boolean[],
	//               $5::smallint[]
	//       ) AS t(m_id, p_id, u, is_p, s_ord)
	BulkCreateProductImages(ctx context.Context, arg BulkCreateProductImagesParams) error
	// 创建审核记录，返回新记录ID
	//
	//  INSERT INTO products.product_audits (product_id,
	//                                       merchant_id,
	//                                       old_status,
	//                                       new_status,
	//                                       reason,
	//                                       operator_id)
	//  VALUES ($1, $2, $3, $4, $5, $6)
	//  RETURNING id, created_at
	CreateAuditRecord(ctx context.Context, arg CreateAuditRecordParams) (CreateAuditRecordRow, error)
	// 所有分片表必须：
	// 1. 包含分片键列（merchant_id）
	// 2. 主键必须包含分片键
	// 3. 外键约束需要特殊处理（Citus 不支持跨节点外键）
	// 创建商品主记录，返回生成的ID
	// merchant_id 作为分片键，必须提供
	//
	//
	//  INSERT INTO products.products (name,
	//                                 description,
	//                                 price,
	//                                 status,
	//                                 merchant_id)
	//  VALUES ($1, $2, $3, $4, $5)
	//  RETURNING id, created_at, updated_at
	CreateProduct(ctx context.Context, arg CreateProductParams) (CreateProductRow, error)
	//CreateProductImages
	//
	//  INSERT INTO products.product_images (merchant_id, -- 新增分片键
	//                                       product_id,
	//                                       url,
	//                                       is_primary,
	//                                       sort_order)
	//  VALUES ($1, $2, $3, $4, $5)
	CreateProductImages(ctx context.Context, arg []CreateProductImagesParams) (int64, error)
	// 获取最新审核记录
	//
	//  INSERT INTO products.product_audits (merchant_id, -- 新增分片键
	//                                       product_id,
	//                                       old_status,
	//                                       new_status,
	//                                       reason,
	//                                       operator_id)
	//  VALUES ($1, $2, $3, $4, $5, $6)
	//  RETURNING id, created_at
	GetLatestAudit(ctx context.Context, arg GetLatestAuditParams) (GetLatestAuditRow, error)
	// 乐观锁版本控制
	// 获取商品详情，包含软删除检查
	//
	//
	//  SELECT id,
	//         name,
	//         description,
	//         price,
	//         status,
	//         merchant_id,
	//         created_at,
	//         updated_at
	//  FROM products.products
	//  WHERE id = $1
	//    AND merchant_id = $2
	//    AND deleted_at IS NULL
	GetProduct(ctx context.Context, arg GetProductParams) (GetProductRow, error)
	// 获取商品图片列表，按排序顺序返回
	//
	//  SELECT id, merchant_id, product_id, url, is_primary, sort_order, created_at
	//  FROM products.product_images
	//  WHERE merchant_id = $1
	//    AND product_id = $2 -- 查询必须包含分片键
	//  ORDER BY sort_order
	GetProductImages(ctx context.Context, arg GetProductImagesParams) ([]ProductsProductImages, error)
	// 软删除商品，设置删除时间戳
	//
	//  UPDATE products.products
	//  SET deleted_at = NOW()
	//  WHERE merchant_id = $1
	//    AND id = $2
	//  RETURNING id, merchant_id, name, description, price, status, current_audit_id, category_id, created_at, updated_at, deleted_at
	SoftDeleteProduct(ctx context.Context, arg SoftDeleteProductParams) (ProductsProducts, error)
	// 更新商品基础信息，使用乐观锁控制并发
	//
	//  UPDATE products.products
	//  SET name        = $2,
	//      description = $3,
	//      price       = $4,
	//      status      = $5,
	//      updated_at  = NOW()
	//  WHERE id = $1
	//    AND merchant_id = $6
	//    AND updated_at = $7
	UpdateProduct(ctx context.Context, arg UpdateProductParams) error
	// 更新商品状态并记录当前审核ID
	//
	//  UPDATE products.products
	//  SET status           = $2,
	//      current_audit_id = $3,
	//      updated_at       = NOW()
	//  WHERE id = $1
	//    AND merchant_id = $4
	UpdateProductStatus(ctx context.Context, arg UpdateProductStatusParams) error
}

var _ Querier = (*Queries)(nil)
