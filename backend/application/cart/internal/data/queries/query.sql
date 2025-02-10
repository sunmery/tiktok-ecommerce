-- name: AddItem :one
INSERT INTO cart_schema.carts (user_id, items, updated_at) 
VALUES ($1, $2::jsonb, $3)
RETURNING *;

-- name: GetCart :one
SELECT items
FROM cart_schema.carts
WHERE user_id = $1
ORDER BY updated_at DESC;

-- name: EmptyCart :exec
DELETE 
FROM cart_schema.carts 
WHERE user_id = $1
RETURNING *;

-- name: UpdateItem :one
UPDATE cart_schema.carts
SET items = CASE
    -- 如果存在 product_id 为 $1 的商品，更新它的 quantity
    WHEN items @> jsonb_build_array(jsonb_build_object('product_id', $1::text)) THEN
        jsonb_set(
            items,
            '{0,quantity}',  -- 这里的路径为数组索引形式，假设 product_id 对应数组的第一个位置
            $2::jsonb -- 更新 quantity 字段为 $2
        )
    -- 如果不存在，则向 items 数组中添加新的商品
    ELSE
        items || jsonb_build_array(
            jsonb_build_object(
                'product_id', $1::text,
                'quantity', $2::int
            )
        )
    END
WHERE user_id = $3
RETURNING *;

-- name: RemoveItem :one
UPDATE cart_schema.carts
SET items = (
    SELECT jsonb_agg(item)
    FROM jsonb_array_elements(items) AS item
    WHERE item->>'product_id' != $1::text
)
WHERE user_id = $2
RETURNING *;
