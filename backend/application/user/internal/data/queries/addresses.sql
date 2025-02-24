-- name: CreatAddress :one
INSERT INTO users.addresses(user_id, street_address, city, state, country, zip_code)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAddresses :many
SELECT *
FROM users.addresses
WHERE user_id = @user_id;

-- name: UpdateAddress :one
UPDATE users.addresses
SET street_address = coalesce(sqlc.narg(street_address), street_address),
    city           = coalesce(sqlc.narg(city), city),
    state          = coalesce(sqlc.narg(state), state),
    country        = coalesce(sqlc.narg(country), country),
    zip_code       = coalesce(sqlc.narg(zip_code), zip_code)
WHERE id = @id
  AND user_id = @user_id

RETURNING *;

-- name: DeleteAddress :one
DELETE
FROM users.addresses
WHERE id = @id
  AND user_id = @user_id
RETURNING *;
