-- name: CreateOrder :one
INSERT INTO orders.orders (owner, name, email, street_address, city, state, country, zip_code, currency, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, now())
RETURNING id;


-- name: ListOrders :many
SELECT * FROM orders.orders
WHERE owner = @owner AND name = @name
ORDER BY updated_at;

-- name: MarkOrderPaid :one
UPDATE orders.orders
SET status = 'paid', updated_at = now()
WHERE id = @id AND owner = @owner AND name = @name
RETURNING *;

-- name: CreateOrderItems :one
INSERT INTO orders.order_items(order_id, product_id, name, quantity, price)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListOrderItems :many
SELECT * FROM orders.order_items
WHERE order_id = @order_id;
