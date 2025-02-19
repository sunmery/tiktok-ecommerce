-- 所有分片表必须：
-- 1. 包含分片键列（merchant_id）
-- 2. 主键必须包含分片键
-- 3. 外键约束需要特殊处理（Citus 不支持跨节点外键）

-- 创建商品主记录，返回生成的ID
-- merchant_id 作为分片键，必须提供
-- name: CreateProduct :one
INSERT INTO products.products (name,
                               description,
                               price,
                               status,
                               merchant_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at;

-- name: UpdateProduct :exec
-- 更新商品基础信息，使用乐观锁控制并发
UPDATE products.products
SET name        = $2,
    description = $3,
    price       = $4,
    status      = $5,
    updated_at  = NOW()
WHERE id = $1
  AND merchant_id = $6
  AND updated_at = $7;
-- 乐观锁版本控制

-- name: GetProduct :one
-- 获取商品详情，包含软删除检查
SELECT id,
       name,
       description,
       price,
       status,
       merchant_id,
       created_at,
       updated_at
FROM products.products
WHERE id = $1
  AND merchant_id = $2
  AND deleted_at IS NULL;

-- name: SoftDeleteProduct :one
-- 软删除商品，设置删除时间戳
UPDATE products.products
SET deleted_at = NOW()
WHERE merchant_id = $1
  AND id = $2
RETURNING *;

-- name: CreateProductImages :copyfrom
INSERT INTO products.product_images (merchant_id, -- 新增分片键
                                     product_id,
                                     url,
                                     is_primary,
                                     sort_order)
VALUES ($1, $2, $3, $4, $5);

-- name: GetProductImages :many
-- 获取商品图片列表，按排序顺序返回
SELECT *
FROM products.product_images
WHERE merchant_id = $1
  AND product_id = $2 -- 查询必须包含分片键
ORDER BY sort_order;

-- 批量插入图片
-- name: BulkCreateProductImages :exec
INSERT INTO products.product_images
    (merchant_id, product_id, url, is_primary, sort_order)
SELECT m_id, p_id, u, is_p, s_ord
FROM unnest(
             @merchant_ids::bigint[],
             @product_ids::bigint[],
             @urls::text[],
             @is_primary::boolean[],
             @sort_orders::smallint[]
     ) AS t(m_id, p_id, u, is_p, s_ord);

-- name: CreateAuditRecord :one
-- 创建审核记录，返回新记录ID
INSERT INTO products.product_audits (product_id,
                                     merchant_id,
                                     old_status,
                                     new_status,
                                     reason,
                                     operator_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at;


-- name: UpdateProductStatus :exec
-- 更新商品状态并记录当前审核ID
UPDATE products.products
SET status           = $2,
    current_audit_id = $3,
    updated_at       = NOW()
WHERE id = $1
  AND merchant_id = $4;

-- name: GetLatestAudit :one
-- 获取最新审核记录
INSERT INTO products.product_audits (merchant_id, -- 新增分片键
                                     product_id,
                                     old_status,
                                     new_status,
                                     reason,
                                     operator_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at;
