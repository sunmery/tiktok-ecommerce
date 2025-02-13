-- name: CreateProduct :one
-- 创建商品
INSERT INTO products.products(name,
                              description,
                              picture,
                              price,
                              total_stock)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: CreateCategories :one
-- 创建分类数据
INSERT INTO products.categories (name, parent_id)
VALUES ($1, $2)
RETURNING *;

-- name: CreateProductCategories :one
-- 关联商品与分类
-- 将商品1关联到分类2（Smartphones）
INSERT INTO products.product_categories (product_id, category_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetProductCategories :many
-- 查询某分类下的所有商品
SELECT p.*
FROM products.products p
         JOIN products.product_categories pc ON p.id = pc.product_id
WHERE pc.category_id = $1;

-- name: UpdateProductsReservedStock :one
-- 预留库存（下单时）
UPDATE products.products
SET reserved_stock = reserved_stock + 2
WHERE id = 1
RETURNING *;

-- name: CreateProductInventoryHistory :one
-- 记录库存变更
INSERT INTO products.inventory_history (product_id,
                                        old_stock,
                                        new_stock,
                                        change_reason)
VALUES ($1,
        (SELECT total_stock FROM products.products WHERE id = $1),
        (SELECT total_stock FROM products.products WHERE id = $1) - $2,
        'ORDER_RESERVED')
RETURNING *;

-- -- name: CreateAuditLog :one
-- INSERT INTO products.audit_log (action, product_id, owner, name)
-- VALUES ($1, $2, $3, $4)
-- RETURNING *;

-- name: ListProducts :many
SELECT *
FROM products.products
ORDER BY id
OFFSET $1 LIMIT $2;

-- name: GetProduct :one
SELECT *
FROM products.products
WHERE id = $1
LIMIT 1;

-- name: SearchProducts :many
SELECT *
FROM products.products
WHERE name ILIKE '%' || $1 || '%';
