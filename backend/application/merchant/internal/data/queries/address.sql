-- addresses.sql
-- 创建商家地址
-- name: CreateAddress :one
INSERT INTO merchant.addresses (id,
                                 merchant_id,
                                 address_type,
                                 contact_person,
                                 contact_phone,
                                 street_address,
                                 city,
                                 state,
                                 country,
                                 zip_code,
                                 is_default,
                                 remarks)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- 批量创建地址（需要服务层处理）
-- name: BatchCreateAddresses :many
INSERT INTO merchant.addresses (id,
                                 merchant_id,
                                 address_type,
                                 contact_person,
                                 contact_phone,
                                 street_address,
                                 city,
                                 state,
                                 country,
                                 zip_code,
                                 is_default,
                                 remarks)
VALUES (UNNEST(@id::bigint[]),
        UNNEST(@merchant_id::uuid[]),
        UNNEST(@address_type::varchar[]),
        UNNEST(@contact_person::varchar[]),
        UNNEST(@contact_phone::varchar[]),
        UNNEST(@street_address::text[]),
        UNNEST(@city::varchar[]),
        UNNEST(@state::varchar[]),
        UNNEST(@country::varchar[]),
        UNNEST(@zip_code::varchar[]),
        UNNEST(@is_default::boolean[]),
        UNNEST(@remarks::text[]))
RETURNING *;

-- 更新地址（带动态字段更新）
-- name: UpdateAddress :one
UPDATE merchant.addresses
SET address_type   = COALESCE($1, address_type),
    contact_person = COALESCE($2, contact_person),
    contact_phone  = COALESCE($3, contact_phone),
    street_address = COALESCE($4, street_address),
    city           = COALESCE($5, city),
    state          = COALESCE($6, state),
    country        = COALESCE($7, country),
    zip_code       = COALESCE($8, zip_code),
    is_default     = COALESCE($9, is_default),
    remarks        = COALESCE($10, remarks),
    updated_at     = NOW()
WHERE id = $11
  AND merchant_id = $12
RETURNING *;

-- 删除地址
-- name: DeleteAddress :exec
DELETE
FROM merchant.addresses
WHERE id = $1
  AND merchant_id = $2;

-- 获取地址详情
-- name: GetAddress :one
SELECT *
FROM merchant.addresses
WHERE id = $1
  AND merchant_id = $2;

-- 地址列表（带分页和过滤）
-- name: ListAddresses :many
SELECT *
FROM merchant.addresses
WHERE merchant_id = $1
  AND (address_type = $2 OR $2 IS NULL)
  AND (is_default = $3 OR $3 IS NULL)
ORDER BY id
LIMIT $4 OFFSET $5;

-- 设置默认地址（带事务处理）
-- name: SetDefaultAddress :one
WITH update_all AS (
    UPDATE merchant.addresses
        SET is_default = false
        WHERE merchant_id = @merchant_id
            AND address_type = (SELECT address_type FROM merchant.addresses WHERE id = $2)
            AND id != @id)
UPDATE merchant.addresses
SET is_default = true
WHERE id = @id
RETURNING *;

-- 智能获取发货地址
-- name: GetShippingAddress :one
SELECT *
FROM merchant.addresses
WHERE merchant_id = $1
  AND address_type = 'WAREHOUSE'
ORDER BY is_default DESC, created_at DESC
LIMIT 1;