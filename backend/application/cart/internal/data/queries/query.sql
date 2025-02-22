-- name: UpsertItem :one
WITH cart_id_cte AS (
    SELECT c.cart_id
    FROM cart_schema.cart AS c
    WHERE c.user_id = $1 AND c.cart_name = $2
    LIMIT 1
),
insert_cart AS (
    INSERT INTO cart_schema.cart (user_id, cart_name)
    SELECT $1, $2
    WHERE NOT EXISTS (SELECT 1 FROM cart_id_cte)
    RETURNING cart_id
)
INSERT INTO cart_schema.cart_items (cart_id, merchant_id, product_id, quantity, created_at, updated_at)
VALUES (
    COALESCE((SELECT cart_id FROM cart_id_cte), (SELECT cart_id FROM insert_cart)),  -- 获取或创建购物车ID
    $3,   -- 商家ID
    $4,   -- 商品ID
    $5,   -- 商品数量
    CURRENT_TIMESTAMP,  -- 创建时间
    CURRENT_TIMESTAMP   -- 更新时间
)
ON CONFLICT (cart_id, merchant_id, product_id)  -- 如果购物车ID、商家ID和商品ID组合重复
DO UPDATE SET 
    quantity = EXCLUDED.quantity,  -- 更新商品数量
    updated_at = CURRENT_TIMESTAMP  -- 更新时间
RETURNING *;

-- name: GetCart :many
SELECT ci.merchant_id, ci.product_id, ci.quantity, ci.selected 
FROM cart_schema.cart_items AS ci
WHERE ci.cart_id = 
    (SELECT c.cart_id
     FROM cart_schema.cart AS c
     WHERE c.user_id = $1 AND c.cart_name = $2 LIMIT 1);  -- 获取用户的购物车ID

-- name: RemoveCartItem :one
DELETE FROM cart_schema.cart_items AS ci
WHERE ci.cart_id = 
    (SELECT c.cart_id
     FROM cart_schema.cart AS c
     WHERE c.user_id = $1 AND c.cart_name = $2 LIMIT 1)  -- 获取用户的购物车ID
    AND ci.merchant_id = $3  -- 商家ID
    AND ci.product_id = $4  -- 删除指定商品ID
RETURNING *;

-- name: EmptyCart :exec
DELETE FROM cart_schema.cart_items AS ci
WHERE ci.cart_id = 
    (SELECT c.cart_id
     FROM cart_schema.cart AS c
     WHERE c.user_id = $1 AND c.cart_name = $2);  -- 获取用户的购物车ID

-- name: CheckCartItem :exec
UPDATE cart_schema.cart_items AS ci
SET selected = TRUE
WHERE ci.cart_id = 
    (SELECT c.cart_id
     FROM cart_schema.cart AS c
     WHERE c.user_id = $1 AND c.cart_name = $2 LIMIT 1) 
    AND ci.merchant_id = $3  -- 商家ID
    AND ci.product_id = $4;

-- name: CreateCart :one
INSERT INTO cart_schema.cart (user_id, cart_name)
VALUES ($1, $2)
RETURNING *;

-- name: CreateOrder :many
SELECT ci.merchant_id, ci.product_id, ci.quantity, ci.selected
FROM cart_schema.cart_items AS ci
WHERE ci.cart_id = 
    (SELECT c.cart_id
     FROM cart_schema.cart AS c
     WHERE c.user_id = $1 AND c.cart_name = $2 LIMIT 1) 
    AND ci.selected = TRUE;

-- name: ListCarts :many
SELECT c.cart_id, c.cart_name
FROM cart_schema.cart AS c
WHERE c.user_id = $1;

-- name: UncheckCartItem :exec
UPDATE cart_schema.cart_items AS ci
SET selected = FALSE
WHERE ci.cart_id = 
    (SELECT c.cart_id
     FROM cart_schema.cart AS c
     WHERE c.user_id = $1 AND c.cart_name = $2 LIMIT 1) 
    AND ci.merchant_id = $3  -- 商家ID
    AND ci.product_id = $4;