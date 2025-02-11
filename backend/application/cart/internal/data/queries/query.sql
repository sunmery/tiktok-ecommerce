-- name: UpsertItem :one
INSERT INTO cart_schema.cart_items (cart_id, product_id, quantity, created_at, updated_at)
VALUES (
    (SELECT c.cart_id
     FROM cart_schema.cart AS c
     WHERE c.user_id = $1 LIMIT 1),  -- 获取用户的购物车ID
    $2,   -- 商品ID
    $3,   -- 商品数量
    CURRENT_TIMESTAMP,  -- 创建时间
    CURRENT_TIMESTAMP   -- 更新时间
)
ON CONFLICT (cart_id, product_id)  -- 如果购物车ID和商品ID组合重复
DO UPDATE SET 
    quantity = cart_schema.cart_items.quantity + EXCLUDED.quantity,  -- 更新商品数量
    updated_at = CURRENT_TIMESTAMP  -- 更新时间
RETURNING *;

-- name: GetCart :many
SELECT ci.cart_item_id, ci.quantity 
FROM cart_schema.cart_items AS ci
WHERE ci.cart_id = 
    (SELECT c.cart_id
     FROM cart_schema.cart AS c
     WHERE c.user_id = $1);  -- 获取用户的购物车ID

-- name: RemoveCartItem :one
DELETE FROM cart_schema.cart_items AS ci
WHERE ci.cart_id = 
    (SELECT c.cart_id
     FROM cart_schema.cart AS c
     WHERE c.user_id = $1)  -- 获取用户的购物车ID
    AND ci.product_id = $2  -- 删除指定商品ID
RETURNING *;

-- name: EmptyCart :many
DELETE FROM cart_schema.cart_items AS ci
WHERE ci.cart_id = 
    (SELECT c.cart_id
     FROM cart_schema.cart AS c
     WHERE c.user_id = $1)  -- 获取用户的购物车ID
RETURNING *;