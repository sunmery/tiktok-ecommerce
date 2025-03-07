-- name: InsertCreditCard :exec
INSERT INTO users.credit_cards (user_id, currency,number, cvv, exp_year, exp_month, owner, name, type, brand, country)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);

-- name: DeleteCreditCard :exec
DELETE
FROM users.credit_cards
WHERE id = $1;

-- name: ListCreditCards :many
SELECT *
FROM users.credit_cards
WHERE user_id = $1;
