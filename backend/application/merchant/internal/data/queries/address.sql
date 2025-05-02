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

-- 获取地址详情
-- name: GetAddress :one
SELECT *
FROM merchant.addresses
WHERE id = $1
  AND merchant_id = $2;

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
        UNNEST(@address_type::merchant.address_type[]),
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

-- 地址列表（带分页和过滤）
-- name: ListFilterAddresses :many
SELECT *
FROM merchant.addresses
WHERE merchant_id = $1
  AND (address_type = $2 OR $2 IS NULL)
ORDER BY id
LIMIT $3 OFFSET $4;

-- 根据地址类型查询默认地址
-- name: GetDefaultAddress :one
SELECT *
FROM merchant.addresses
WHERE merchant_id = @merchant_id
  AND address_type = @address_type
  AND is_default = true;

-- 获取全部地址的默认值列表
-- name: GetDefaultAddresses :many
SELECT *
FROM merchant.addresses
WHERE merchant_id = @merchant_id
  AND is_default = true;

-- 查询全部地址（带分页）
-- name: ListAddresses :many
SELECT *
FROM merchant.addresses
WHERE merchant_id = $1
ORDER BY id
LIMIT $2 OFFSET $3;

-- 设置默认地址（带事务处理）
-- name: SetDefaultAddress :one
WITH get_address_type AS (SELECT address_type
                          FROM merchant.addresses
                          WHERE id = @id
                            AND merchant_id = @merchant_id),
     update_old_default AS (
         UPDATE merchant.addresses
             SET is_default = false
             WHERE merchant_id = @merchant_id
                 AND address_type = (SELECT address_type FROM get_address_type)
                 AND is_default = true)
UPDATE merchant.addresses
SET is_default = true
WHERE id = @id
  AND merchant_id = @merchant_id
RETURNING *;
