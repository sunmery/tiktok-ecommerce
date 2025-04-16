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
                               merchant_id,
                               category_id
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at, updated_at;

-- 获取商品详情，包含软删除检查
-- name: GetProduct :one
SELECT p.id,
       p.name,
       p.description,
       p.price,
       p.status,
       p.merchant_id,
       p.category_id,
       p.created_at,
       p.updated_at,
       i.stock,
       (SELECT jsonb_agg(jsonb_build_object(
               'url', pi.url,
               'is_primary', pi.is_primary,
               'sort_order', pi.sort_order
                         ))
        FROM products.product_images pi
        WHERE pi.product_id = p.id
          AND pi.merchant_id = p.merchant_id) AS images,
       pa.attributes,
       (SELECT jsonb_build_object(
                       'id', a.id,
                       'old_status', a.old_status,
                       'new_status', a.new_status,
                       'reason', a.reason,
                       'created_at', a.created_at
               )
        FROM products.product_audits a
        WHERE a.product_id = p.id
          AND a.merchant_id = p.merchant_id
        ORDER BY a.created_at DESC
        LIMIT 1)                              AS latest_audit
FROM products.products p
         INNER JOIN products.inventory i
                    ON p.id = i.product_id AND p.merchant_id = i.merchant_id
         LEFT JOIN products.product_attributes pa
                   ON p.id = pa.product_id AND p.merchant_id = pa.merchant_id
WHERE p.id = $1
  AND p.merchant_id = $2
  AND p.deleted_at IS NULL;

-- name: GetProductsBatch :many
SELECT p.id,
       p.name,
       p.description,
       p.price,
       p.status,
       p.merchant_id,
       p.category_id,
       p.created_at,
       p.updated_at,
       i.stock,
       (SELECT jsonb_agg(jsonb_build_object(
               'url', pi.url,
               'is_primary', pi.is_primary,
               'sort_order', pi.sort_order
                         ))
        FROM products.product_images pi
        WHERE pi.product_id = p.id
          AND pi.merchant_id = p.merchant_id) AS images,
       pa.attributes,
       (SELECT jsonb_build_object(
                       'id', a.id,
                       'old_status', a.old_status,
                       'new_status', a.new_status,
                       'reason', a.reason,
                       'created_at', a.created_at
               )
        FROM products.product_audits a
        WHERE a.product_id = p.id
          AND a.merchant_id = p.merchant_id
        ORDER BY a.created_at DESC
        LIMIT 1)                              AS latest_audit
FROM products.products p
         INNER JOIN products.inventory i
                    ON p.id = i.product_id AND p.merchant_id = i.merchant_id
         LEFT JOIN products.product_attributes pa
                   ON p.id = pa.product_id AND p.merchant_id = pa.merchant_id
WHERE p.id = ANY(@product_ids::UUID[])
  AND p.merchant_id = ANY(@merchant_ids::UUID[])
  AND p.deleted_at IS NULL;



-- 商品搜索查询
-- name: SearchFullProductsByName :many
SELECT p.id,
       p.name,
       p.description,
       p.price,
       p.status,
       p.merchant_id,
       p.created_at,
       p.updated_at,
       i.stock,
       (SELECT jsonb_agg(jsonb_build_object(
               'url', pi.url,
               'is_primary', pi.is_primary,
               'sort_order', pi.sort_order
                         ))
        FROM products.product_images pi
        WHERE pi.product_id = p.id
          AND pi.merchant_id = p.merchant_id) AS images,
       pa.attributes
FROM products.products p
         INNER JOIN products.inventory i
                    ON p.id = i.product_id AND p.merchant_id = i.merchant_id
         LEFT JOIN products.product_attributes pa
                   ON p.id = pa.product_id AND p.merchant_id = pa.merchant_id
WHERE p.name ILIKE '%' || @name || '%'
  AND p.deleted_at IS NULL
ORDER BY ts_rank(to_tsvector('simple', p.name), plainto_tsquery('simple', @query)) DESC,
         p.created_at DESC
LIMIT @page OFFSET @page_size;

-- 软删除商品，设置删除时间戳
-- name: SoftDeleteProduct :one
UPDATE products.products
SET deleted_at = NOW(),
    status     = $3
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
             @merchant_ids::uuid[],
             @product_ids::uuid[],
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

-- 实现随机商品列表查询（使用PostgreSQL的TABLESAMPLE优化性能）,数据量太少不显示
-- name: ListRandomProducts :many
-- SELECT
--     p.id,
--     p.merchant_id,
--     p.name,
--     p.description,
--     p.price,
--     p.status,
--     p.category_id,
--     p.created_at,
--     p.updated_at,
--     -- 图片信息
--     (
--         SELECT jsonb_agg(jsonb_build_object(
--                 'url', pi.url,
--                 'is_primary', pi.is_primary,
--                 'sort_order', pi.sort_order
--                          ))
--         FROM products.product_images pi
--         WHERE pi.product_id = p.id AND pi.merchant_id = p.merchant_id
--     ) AS images,
--     -- 属性信息
--     pa.attributes
-- FROM products.products p
--          TABLESAMPLE BERNOULLI (0.1) REPEATABLE (123)
--          LEFT JOIN products.product_attributes pa
--                    ON p.id = pa.product_id AND p.merchant_id = pa.merchant_id
-- WHERE p.status = $1 AND p.deleted_at IS NULL
-- ORDER BY random()
--     LIMIT $2 OFFSET $3;

-- name: ListRandomProducts :many
SELECT p.id,
       p.merchant_id,
       p.name,
       p.description,
       p.price,
       p.status,
       p.category_id,
       p.created_at,
       p.updated_at,
       i.stock,
       -- 图片信息
       (SELECT jsonb_agg(jsonb_build_object(
               'url', pi.url,
               'is_primary', pi.is_primary,
               'sort_order', pi.sort_order
                         ))
        FROM products.product_images pi
        WHERE pi.product_id = p.id
          AND pi.merchant_id = p.merchant_id) AS images,
       -- 属性信息
       pa.attributes
FROM products.products p
         INNER JOIN products.inventory i
                    ON p.id = i.product_id AND p.merchant_id = i.merchant_id
         LEFT JOIN products.product_attributes pa
                   ON p.id = pa.product_id AND p.merchant_id = pa.merchant_id
WHERE p.status = $1
  AND p.deleted_at IS NULL
ORDER BY random()
LIMIT $2 OFFSET $3;

-- 根据分类获取商品列表
-- name: GetCategoryProducts :many
WITH filtered_products AS (SELECT p.id,
                                  p.merchant_id,
                                  p.name,
                                  p.description,
                                  p.price,
                                  p.status,
                                  p.category_id,
                                  p.created_at,
                                  p.updated_at
                           FROM products.products p
                           WHERE p.category_id = $1 -- 指定分类id
                             AND p.status = $2      -- 商品状态机
                             AND p.deleted_at IS NULL),
     product_images_agg AS (SELECT pi.product_id,
                                   jsonb_agg(
                                           jsonb_build_object(
                                                   'id', pi.id,
                                                   'url', pi.url,
                                                   'is_primary', pi.is_primary,
                                                   'sort_order', pi.sort_order
                                           )
                                   ) AS images
                            FROM products.product_images pi
                                     INNER JOIN filtered_products fp
                                                ON pi.product_id = fp.id AND pi.merchant_id = fp.merchant_id
                            GROUP BY pi.product_id),
     product_attributes_agg AS (SELECT pa.product_id,
                                       pa.attributes
                                FROM products.product_attributes pa
                                         INNER JOIN filtered_products fp
                                                    ON pa.product_id = fp.id AND pa.merchant_id = fp.merchant_id)
SELECT fp.id,
       fp.merchant_id,
       fp.name,
       fp.description,
       fp.price,
       fp.status,
       fp.category_id,
       fp.created_at,
       fp.updated_at,
       COALESCE(pia.images, '[]'::jsonb)     AS images,
       COALESCE(paa.attributes, '{}'::jsonb) AS attributes
FROM filtered_products fp
         LEFT JOIN product_images_agg pia
                   ON fp.id = pia.product_id
         LEFT JOIN product_attributes_agg paa
                   ON fp.id = paa.product_id
ORDER BY fp.created_at DESC
LIMIT $3 OFFSET $4;

-- 根据分类及其所有子分类获取商品列表
-- name: GetCategoryWithChildrenProducts :many
WITH RECURSIVE category_hierarchy AS (
    -- 基础情况：指定的分类
    SELECT id, parent_id
    FROM categories.categories
    WHERE id = @id
    
    UNION ALL
    
    -- 递归情况：所有子分类
    SELECT c.id, c.parent_id
    FROM categories.categories c
    JOIN category_hierarchy ch ON c.parent_id = ch.id
),
filtered_products AS (
    SELECT p.id,
           p.merchant_id,
           p.name,
           p.description,
           p.price,
           p.status,
           p.category_id,
           p.created_at,
           p.updated_at
    FROM products.products p
    JOIN category_hierarchy ch ON p.category_id = ch.id
    WHERE p.status = @status      -- 商品状态机
      AND p.deleted_at IS NULL
),
product_images_agg AS (
    SELECT pi.product_id,
           jsonb_agg(
               jsonb_build_object(
                   'id', pi.id,
                   'url', pi.url,
                   'is_primary', pi.is_primary,
                   'sort_order', pi.sort_order
               )
           ) AS images
    FROM products.product_images pi
    INNER JOIN filtered_products fp ON pi.product_id = fp.id AND pi.merchant_id = fp.merchant_id
    GROUP BY pi.product_id
),
product_attributes_agg AS (
    SELECT pa.product_id,
           pa.attributes
    FROM products.product_attributes pa
    INNER JOIN filtered_products fp ON pa.product_id = fp.id AND pa.merchant_id = fp.merchant_id
),
inventory_agg AS (
    SELECT i.product_id,
           i.stock
    FROM products.inventory i
    INNER JOIN filtered_products fp ON i.product_id = fp.id AND i.merchant_id = fp.merchant_id
)
SELECT fp.id,
       fp.merchant_id,
       fp.name,
       fp.description,
       fp.price,
       fp.status,
       fp.category_id,
       fp.created_at,
       fp.updated_at,
       COALESCE(ia.stock, 0) AS stock,
       COALESCE(pia.images, '[]'::jsonb) AS images,
       COALESCE(paa.attributes, '{}'::jsonb) AS attributes
FROM filtered_products fp
LEFT JOIN product_images_agg pia ON fp.id = pia.product_id
LEFT JOIN product_attributes_agg paa ON fp.id = paa.product_id
LEFT JOIN inventory_agg ia ON fp.id = ia.product_id
ORDER BY fp.created_at DESC
LIMIT @page_size OFFSET @page;

-- 分类批量查询（使用GIN索引优化数组查询）
-- name: ListProductsByCategory :many
SELECT p.id,
       p.merchant_id,
       p.name,
       p.description,
       p.price,
       p.status,
       p.category_id,
       p.created_at,
       p.updated_at,
       i.stock,
       (SELECT jsonb_agg(jsonb_build_object(
               'url', pi.url,
               'is_primary', pi.is_primary,
               'sort_order', pi.sort_order
                         ))
        FROM products.product_images pi
        WHERE pi.product_id = p.id
          AND pi.merchant_id = p.merchant_id) AS images,
       pa.attributes
FROM products.products p
         INNER JOIN products.inventory i
                    ON p.id = i.product_id AND p.merchant_id = i.merchant_id
         LEFT JOIN products.product_attributes pa
                   ON p.id = pa.product_id AND p.merchant_id = pa.merchant_id
WHERE p.category_id = ANY ($1::bigint[])
  AND p.status = $2
  AND p.deleted_at IS NULL
ORDER BY p.created_at DESC
LIMIT $3 OFFSET $4;

-- name: CreateProductAttribute :exec
INSERT INTO products.product_attributes
    (merchant_id, product_id, attributes)
VALUES ($1, $2, $3)
RETURNING created_at, updated_at;

-- name: GetProductAttribute :one
SELECT merchant_id,
       product_id,
       attributes,
       created_at,
       updated_at
FROM products.product_attributes
WHERE merchant_id = $1
  AND product_id = $2;

-- name: UpdateProductAttribute :exec
UPDATE products.product_attributes
SET attributes = $3,
    updated_at = NOW()
WHERE merchant_id = $1
  AND product_id = $2
RETURNING updated_at;

-- name: DeleteProductAttribute :exec
DELETE
FROM products.product_attributes
WHERE merchant_id = $1
  AND product_id = $2;
