-- name: CreateProduct :one
-- 创建商品
INSERT INTO products.products(name,
                              description,
                              picture,
                              price,
                              category_id,
                              total_stock)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: CreateAuditLog :one
-- 创建审计日志
INSERT INTO products.inventory_history(change_reason, product_id, new_stock, old_stock, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateAuditLog :one
-- 更新审计日志
INSERT INTO products.inventory_history (
    product_id, 
    change_reason, 
    new_stock, 
    user_id,
    old_stock
)
VALUES (
    $1,  -- product_id
    $2,  -- change_reason
    $3,  -- new_stock
    $4, -- user_id
    (SELECT total_stock FROM products.products WHERE id = $1)  -- old_stock
)
RETURNING *;

-- name: CreateCategories :one
-- 创建分类数据
INSERT INTO products.categories (name, parent_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetCategories :many
-- 查询所有分类
SELECT *
FROM products.categories;


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


-- name: ListProducts :many

SELECT *
FROM products.products
WHERE ($1 = ANY(category_id))
ORDER BY id
OFFSET $2 LIMIT $3;

-- name: GetProduct :one
SELECT *
FROM products.products
WHERE id = $1
LIMIT 1;

-- name: SearchProducts :many
SELECT *
FROM products.products
WHERE name ILIKE '%' || $1 || '%';

-- name: UpdateProduct :one
UPDATE products.products
SET name = $1, description = $2, picture = $3, price = $4, category_Id = $5, total_stock = $6
WHERE id = $7
RETURNING *;

-- name: DeleteProduct :one
DELETE FROM products.products
WHERE id = @id
RETURNING *;

