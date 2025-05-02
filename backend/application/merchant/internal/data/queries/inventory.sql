-- 获取产品库存
-- name: GetProductStock :one
SELECT p.id                                                                     as product_id,
       p.merchant_id,
       i.stock,
       COALESCE(sa.threshold, 0)                                                as alert_threshold,
       CASE WHEN i.stock <= COALESCE(sa.threshold, 10) THEN true ELSE false END as is_low_stock
FROM products.products p
         JOIN products.inventory i
              ON p.id = i.product_id
                  AND p.merchant_id = i.merchant_id
         LEFT JOIN merchant.stock_alerts sa
                   ON p.id = sa.product_id
                       AND p.merchant_id = sa.merchant_id
WHERE p.id = @product_id::uuid
  AND p.merchant_id = @merchant_id::uuid;

-- 更新产品库存
-- name: UpdateProductStock :one
UPDATE products.inventory
SET stock = stock + @stock
WHERE product_id = @product_id::uuid
  AND merchant_id = @merchant_id::uuid
RETURNING product_id, merchant_id, stock;

-- 设置库存警报阈值
-- name: SetStockAlert :one
INSERT INTO merchant.stock_alerts (product_id, merchant_id, threshold)
VALUES (@product_id::uuid, @merchant_id::uuid, @threshold)
ON CONFLICT (product_id, merchant_id) DO UPDATE
    SET threshold  = @threshold,
        updated_at = NOW()
RETURNING id, merchant.stock_alerts.product_id, merchant_id, threshold, created_at, updated_at;

-- 获取库存警报配置
-- name: GetStockAlerts :many
SELECT sa.id,
       sa.product_id,
       sa.merchant_id,
       p.name  as product_name,
       i.stock as current_stock,
       sa.threshold,
       sa.created_at,
       sa.updated_at
FROM merchant.stock_alerts sa
         JOIN products.products p ON sa.product_id = p.id::uuid -- 显式转换
         JOIN products.inventory i ON sa.product_id = i.product_id::uuid
WHERE sa.merchant_id = @merchant_id::uuid -- 强制类型
ORDER BY sa.updated_at DESC
LIMIT @page_size OFFSET @page;

-- 获取库存警报配置总数
-- name: CountStockAlerts :one
SELECT COUNT(*)
FROM merchant.stock_alerts sa
WHERE sa.merchant_id = @merchant_id::uuid;

-- 获取低库存产品列表
-- name: GetLowStockProducts :many
SELECT p.id                       as product_id,
       p.merchant_id,
       p.name                     as product_name,
       i.stock                    as current_stock,
       COALESCE(sa.threshold, 10) as threshold,
       COALESCE(
               (SELECT pi.url
                FROM products.product_images pi
                WHERE pi.product_id = p.id
                  AND pi.merchant_id = p.merchant_id
                  AND pi.is_primary = true
                LIMIT 1),
               (SELECT pi.url
                FROM products.product_images pi
                WHERE pi.product_id = p.id
                  AND pi.merchant_id = p.merchant_id
                LIMIT 1)
       )                          as image_url
FROM products.products p
         JOIN products.inventory i
              ON p.id = i.product_id
                  AND p.merchant_id = i.merchant_id
         LEFT JOIN merchant.stock_alerts sa
                   ON p.id = sa.product_id
                       AND p.merchant_id = sa.merchant_id
WHERE i.stock <= COALESCE(sa.threshold, 10)
  AND p.merchant_id = @merchant_id::UUID
ORDER BY (i.stock * 1.0 / COALESCE(sa.threshold, 10))
LIMIT @page_size OFFSET @page;

-- 记录库存调整
-- name: RecordStockAdjustment :one
INSERT INTO merchant.stock_adjustments (product_id, merchant_id, quantity, reason, operator_id)
VALUES (@product_id::uuid, @merchant_id::uuid, @quantity, @reason, @operator_id::uuid)
RETURNING id, product_id, merchant_id, quantity, reason, operator_id, created_at;

-- name: GetStockAdjustmentHistory :many
SELECT sa.id,
       sa.product_id,
       sa.merchant_id,
       p.name as product_name,
       sa.quantity,
       sa.reason,
       sa.operator_id,
       sa.created_at
FROM merchant.stock_adjustments sa
         JOIN products.products p
              ON sa.merchant_id = p.merchant_id
WHERE sa.product_id = @product_id::uuid
ORDER BY sa.created_at DESC
LIMIT @page_size OFFSET @page;

-- 获取库存调整历史总数
-- name: CountStockAdjustmentHistory :one
SELECT COUNT(*)
FROM merchant.stock_adjustments sa
  WHERE sa.merchant_id = @merchant_id::uuid;

-- 获取低库存产品总数
-- name: CountLowStockProducts :one
SELECT COUNT(*)
FROM products.products p
         JOIN products.inventory i
              ON p.id = i.product_id
                  AND p.merchant_id = i.merchant_id
         LEFT JOIN merchant.stock_alerts sa
                   ON p.id = sa.product_id
                       AND p.merchant_id = sa.merchant_id
WHERE i.stock <= COALESCE(sa.threshold, @threshold)
  AND p.merchant_id = @merchant_id::uuid;
