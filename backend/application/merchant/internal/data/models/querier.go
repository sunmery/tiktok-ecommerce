// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package models

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	// 获取低库存产品总数
	//
	//  SELECT COUNT(*)
	//  FROM products.products p
	//           JOIN products.inventory i
	//                ON p.id = i.product_id
	//                    AND p.merchant_id = i.merchant_id
	//           LEFT JOIN merchant.stock_alerts sa
	//                     ON p.id = sa.product_id
	//                         AND p.merchant_id = sa.merchant_id
	//  WHERE i.stock <= COALESCE(sa.threshold, $1)
	//    AND p.merchant_id = $2::uuid
	CountLowStockProducts(ctx context.Context, arg CountLowStockProductsParams) (*int64, error)
	// 获取库存调整历史总数
	//
	//  SELECT COUNT(*)
	//  FROM merchant.stock_adjustments sa
	//  WHERE sa.product_id = $1::uuid
	//    AND sa.merchant_id = $2::uuid
	CountStockAdjustmentHistory(ctx context.Context, arg CountStockAdjustmentHistoryParams) (int64, error)
	// 获取库存警报配置总数
	//
	//  SELECT COUNT(*)
	//  FROM merchant.stock_alerts sa
	//  WHERE sa.merchant_id = $1::uuid
	CountStockAlerts(ctx context.Context, merchantID uuid.UUID) (int64, error)
	// 获取低库存产品列表
	//
	//  SELECT p.id                       as product_id,
	//         p.merchant_id,
	//         p.name                     as product_name,
	//         i.stock                    as current_stock,
	//         COALESCE(sa.threshold, 10) as threshold,
	//         COALESCE(
	//                 (SELECT pi.url
	//                  FROM products.product_images pi
	//                  WHERE pi.product_id = p.id
	//                    AND pi.merchant_id = p.merchant_id
	//                    AND pi.is_primary = true
	//                  LIMIT 1),
	//                 (SELECT pi.url
	//                  FROM products.product_images pi
	//                  WHERE pi.product_id = p.id
	//                    AND pi.merchant_id = p.merchant_id
	//                  LIMIT 1)
	//         )                          as image_url
	//  FROM products.products p
	//           JOIN products.inventory i
	//                ON p.id = i.product_id
	//                    AND p.merchant_id = i.merchant_id
	//           LEFT JOIN merchant.stock_alerts sa
	//                     ON p.id = sa.product_id
	//                         AND p.merchant_id = sa.merchant_id
	//  WHERE i.stock <= COALESCE(sa.threshold, 10)
	//    AND p.merchant_id = $1::UUID
	//  ORDER BY (i.stock * 1.0 / COALESCE(sa.threshold, 10))
	//  LIMIT $3 OFFSET $2
	GetLowStockProducts(ctx context.Context, arg GetLowStockProductsParams) ([]GetLowStockProductsRow, error)
	//GetMerchantProducts
	//
	//  SELECT p.id,
	//         p.name,
	//         p.description,
	//         p.price,
	//         p.status,
	//         p.merchant_id,
	//         p.category_id,
	//         p.created_at,
	//         p.updated_at,
	//         i.stock,
	//         (SELECT jsonb_agg(jsonb_build_object(
	//                 'url', pi.url,
	//                 'is_primary', pi.is_primary,
	//                 'sort_order', pi.sort_order
	//                           ))
	//          FROM products.product_images pi
	//          WHERE pi.merchant_id = p.merchant_id) AS images,
	//         pa.attributes,
	//         (SELECT jsonb_build_object(
	//                         'id', a.id,
	//                         'old_status', a.old_status,
	//                         'new_status', a.new_status,
	//                         'reason', a.reason,
	//                         'created_at', a.created_at
	//                 )
	//          FROM products.product_audits a
	//          WHERE a.merchant_id = p.merchant_id
	//          ORDER BY a.created_at DESC
	//          LIMIT 1)                              AS latest_audit
	//  FROM products.products p
	//           INNER JOIN products.inventory i
	//                      ON p.id = i.product_id AND p.merchant_id = i.merchant_id
	//           LEFT JOIN products.product_attributes pa
	//                     ON p.id = pa.product_id AND p.merchant_id = pa.merchant_id
	//  WHERE p.merchant_id = $1
	//    AND p.deleted_at IS NULL
	//  LIMIT $3 OFFSET $2
	GetMerchantProducts(ctx context.Context, arg GetMerchantProductsParams) ([]GetMerchantProductsRow, error)
	// 获取产品库存
	//
	//  SELECT p.id                                                                     as product_id,
	//         p.merchant_id,
	//         i.stock,
	//         COALESCE(sa.threshold, 0)                                                as alert_threshold,
	//         CASE WHEN i.stock <= COALESCE(sa.threshold, 10) THEN true ELSE false END as is_low_stock
	//  FROM products.products p
	//           JOIN products.inventory i
	//                ON p.id = i.product_id
	//                    AND p.merchant_id = i.merchant_id
	//           LEFT JOIN merchant.stock_alerts sa
	//                     ON p.id = sa.product_id
	//                         AND p.merchant_id = sa.merchant_id
	//  WHERE p.id = $1::uuid
	//    AND p.merchant_id = $2::uuid
	GetProductStock(ctx context.Context, arg GetProductStockParams) (GetProductStockRow, error)
	// 获取库存调整历史
	// WHERE sa.product_id = @product_id::uuid
	//
	//  SELECT sa.id,
	//         sa.product_id,
	//         sa.merchant_id,
	//         p.name as product_name,
	//         sa.quantity,
	//         sa.reason,
	//         sa.operator_id,
	//         sa.created_at
	//  FROM merchant.stock_adjustments sa
	//           JOIN products.products p
	//                ON sa.product_id = p.id
	//                    AND sa.merchant_id = p.merchant_id
	//  WHERE sa.merchant_id = $1::uuid
	//  ORDER BY sa.created_at DESC
	//  LIMIT $3 OFFSET $2
	GetStockAdjustmentHistory(ctx context.Context, arg GetStockAdjustmentHistoryParams) ([]GetStockAdjustmentHistoryRow, error)
	// 获取库存警报配置
	//
	//  SELECT sa.id,
	//         sa.product_id,
	//         sa.merchant_id,
	//         p.name  as product_name,
	//         i.stock as current_stock,
	//         sa.threshold,
	//         sa.created_at,
	//         sa.updated_at
	//  FROM merchant.stock_alerts sa
	//           JOIN products.products p ON sa.product_id = p.id::uuid -- 显式转换
	//           JOIN products.inventory i ON sa.product_id = i.product_id::uuid
	//  WHERE sa.merchant_id = $1::uuid -- 强制类型
	//  ORDER BY sa.updated_at DESC
	//  LIMIT $3 OFFSET $2
	GetStockAlerts(ctx context.Context, arg GetStockAlertsParams) ([]GetStockAlertsRow, error)
	//ListOrdersByUser
	//
	//  SELECT id,
	//         order_id,
	//         merchant_id,
	//         total_amount,
	//         currency,
	//         status,
	//         items,
	//         created_at,
	//         updated_at
	//  FROM orders.sub_orders
	//  WHERE merchant_id = $1
	//  ORDER BY created_at DESC
	//  LIMIT $3 OFFSET $2
	ListOrdersByUser(ctx context.Context, arg ListOrdersByUserParams) ([]ListOrdersByUserRow, error)
	//QuerySubOrders
	//
	//  SELECT id,
	//         merchant_id,
	//         total_amount,
	//         currency,
	//         status,
	//         items,
	//         created_at,
	//         updated_at
	//  FROM orders.sub_orders
	//  WHERE order_id = $1
	//  ORDER BY created_at
	QuerySubOrders(ctx context.Context, orderID *int64) ([]QuerySubOrdersRow, error)
	// 记录库存调整
	//
	//  INSERT INTO merchant.stock_adjustments (product_id, merchant_id, quantity, reason, operator_id)
	//  VALUES ($1::uuid, $2::uuid, $3, $4, $5::uuid)
	//  RETURNING id, product_id, merchant_id, quantity, reason, operator_id, created_at
	RecordStockAdjustment(ctx context.Context, arg RecordStockAdjustmentParams) (MerchantStockAdjustments, error)
	// 设置库存警报阈值
	//
	//  INSERT INTO merchant.stock_alerts (product_id, merchant_id, threshold)
	//  VALUES ($1::uuid, $2::uuid, $3)
	//  ON CONFLICT (product_id, merchant_id) DO UPDATE
	//      SET threshold  = $3,
	//          updated_at = NOW()
	//  RETURNING id, merchant.stock_alerts.product_id, merchant_id, threshold, created_at, updated_at
	SetStockAlert(ctx context.Context, arg SetStockAlertParams) (MerchantStockAlerts, error)
	//UpdateProduct
	//
	//  WITH update_product AS (
	//      UPDATE products.products
	//          SET name = coalesce($2, name),
	//              description = coalesce($3, description),
	//              price = coalesce($4, price),
	//              updated_at = now()
	//          WHERE id = $5
	//              AND merchant_id = $6
	//          RETURNING merchant_id,id),
	//       update_attr AS (
	//           UPDATE products.product_attributes
	//               SET attributes = $7,
	//                   updated_at = NOW()
	//               WHERE merchant_id = $6
	//                   AND product_id = $5
	//               RETURNING updated_at),
	//       update_image AS (
	//           UPDATE products.product_images
	//               SET url = $8
	//               WHERE merchant_id = $6
	//                   AND product_id = $5)
	//  UPDATE products.inventory pi
	//  SET stock      = $1,
	//      updated_at = now()
	//  FROM update_product
	//  WHERE update_product.merchant_id = pi.merchant_id
	//    AND update_product.id = pi.product_id
	UpdateProduct(ctx context.Context, arg UpdateProductParams) error
	// 更新产品库存
	//
	//  UPDATE products.inventory
	//  SET stock = stock + $1
	//  WHERE product_id = $2::uuid
	//    AND merchant_id = $3::uuid
	//  RETURNING product_id, merchant_id, stock
	UpdateProductStock(ctx context.Context, arg UpdateProductStockParams) (UpdateProductStockRow, error)
}

var _ Querier = (*Queries)(nil)
