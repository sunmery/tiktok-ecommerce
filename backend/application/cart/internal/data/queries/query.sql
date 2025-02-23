-- name: UpsertCartItem :one
WITH
-- 1. 尝试获取现有购物车
existing_cart AS (
    SELECT cart_id
    FROM carts.carts
    WHERE user_id = @user_id
      AND cart_name = 'cart'
),
-- 2. 若不存在则插入新购物车
inserted_cart AS (
    INSERT INTO carts.carts (user_id, cart_name)
        SELECT @user_id, 'cart'
        WHERE NOT EXISTS (SELECT 1 FROM existing_cart) -- 仅在无现有购物车时插入
        ON CONFLICT (user_id, cart_name) DO NOTHING    -- 处理并发插入冲突
        RETURNING cart_id
),
-- 3. 合并现有或新插入的购物车ID
target_cart AS (
    SELECT cart_id FROM existing_cart
    UNION ALL
    SELECT cart_id FROM inserted_cart
)
-- 4. 插入或更新商品项
INSERT INTO carts.cart_items (cart_id, product_id, quantity)
SELECT cart_id, @product_id, @quantity
FROM target_cart
ON CONFLICT (cart_id, product_id)
    DO UPDATE SET
                  quantity = cart_items.quantity + EXCLUDED.quantity,
                  updated_at = now()
RETURNING *;

-- name: GetCart :many
SELECT ci.product_id, ci.quantity
FROM carts.cart_items AS ci
WHERE ci.cart_id =
      (SELECT c.cart_id
       FROM carts.carts AS c
       WHERE c.user_id = $1
         AND c.cart_name = $2
       LIMIT 1);
-- 获取用户的购物车ID


-- name: RemoveCartItem :one
DELETE
FROM carts.cart_items AS ci
WHERE ci.cart_id =
      (SELECT c.cart_id
       FROM carts.carts AS c
       WHERE c.user_id = $1
         AND c.cart_name = $2
       LIMIT 1)          -- 获取用户的购物车ID
  AND ci.product_id = $3 -- 删除指定商品ID
RETURNING *;

-- name: EmptyCart :exec
DELETE
FROM carts.cart_items AS ci
WHERE ci.cart_id =
      (SELECT c.cart_id
       FROM carts.carts AS c
       WHERE c.user_id = $1
         AND c.cart_name = $2); -- 获取用户的购物车ID
