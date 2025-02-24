// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const CreateOrder = `-- name: CreateOrder :one
INSERT INTO orders.orders (user_id, currency, street_address,
                           city, state, country, zip_code, email)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, user_id, currency, street_address, city, state, country, zip_code, email, created_at, updated_at, payment_status
`

type CreateOrderParams struct {
	UserID        uuid.UUID `json:"userID"`
	Currency      string    `json:"currency"`
	StreetAddress string    `json:"streetAddress"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	Country       string    `json:"country"`
	ZipCode       string    `json:"zipCode"`
	Email         string    `json:"email"`
}

// CreateOrder
//
//	INSERT INTO orders.orders (user_id, currency, street_address,
//	                           city, state, country, zip_code, email)
//	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
//	RETURNING id, user_id, currency, street_address, city, state, country, zip_code, email, created_at, updated_at, payment_status
func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (OrdersOrders, error) {
	row := q.db.QueryRow(ctx, CreateOrder,
		arg.UserID,
		arg.Currency,
		arg.StreetAddress,
		arg.City,
		arg.State,
		arg.Country,
		arg.ZipCode,
		arg.Email,
	)
	var i OrdersOrders
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Currency,
		&i.StreetAddress,
		&i.City,
		&i.State,
		&i.Country,
		&i.ZipCode,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PaymentStatus,
	)
	return i, err
}

const CreateSubOrder = `-- name: CreateSubOrder :one
INSERT INTO orders.sub_orders (order_id, merchant_id, total_amount,
                               currency, status, items)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, order_id, merchant_id, total_amount, currency, status, items, created_at, updated_at, payment_status
`

type CreateSubOrderParams struct {
	OrderID     uuid.UUID      `json:"orderID"`
	MerchantID  uuid.UUID      `json:"merchantID"`
	TotalAmount pgtype.Numeric `json:"totalAmount"`
	Currency    string         `json:"currency"`
	Status      string         `json:"status"`
	Items       []byte         `json:"items"`
}

// CreateSubOrder
//
//	INSERT INTO orders.sub_orders (order_id, merchant_id, total_amount,
//	                               currency, status, items)
//	VALUES ($1, $2, $3, $4, $5, $6)
//	RETURNING id, order_id, merchant_id, total_amount, currency, status, items, created_at, updated_at, payment_status
func (q *Queries) CreateSubOrder(ctx context.Context, arg CreateSubOrderParams) (OrdersSubOrders, error) {
	row := q.db.QueryRow(ctx, CreateSubOrder,
		arg.OrderID,
		arg.MerchantID,
		arg.TotalAmount,
		arg.Currency,
		arg.Status,
		arg.Items,
	)
	var i OrdersSubOrders
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.MerchantID,
		&i.TotalAmount,
		&i.Currency,
		&i.Status,
		&i.Items,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PaymentStatus,
	)
	return i, err
}

const GetOrderByID = `-- name: GetOrderByID :one
SELECT id, user_id, currency, street_address, city, state, country, zip_code, email, created_at, updated_at, payment_status
FROM orders.orders
WHERE id = $1
`

// GetOrderByID
//
//	SELECT id, user_id, currency, street_address, city, state, country, zip_code, email, created_at, updated_at, payment_status
//	FROM orders.orders
//	WHERE id = $1
func (q *Queries) GetOrderByID(ctx context.Context, id uuid.UUID) (OrdersOrders, error) {
	row := q.db.QueryRow(ctx, GetOrderByID, id)
	var i OrdersOrders
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Currency,
		&i.StreetAddress,
		&i.City,
		&i.State,
		&i.Country,
		&i.ZipCode,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PaymentStatus,
	)
	return i, err
}

const GetUserOrdersWithSuborders = `-- name: GetUserOrdersWithSuborders :many

SELECT o.id         AS order_id,
       o.currency   AS order_currency,
       o.street_address,
       o.city,
       o.state,
       o.country,
       o.zip_code,
       o.email,
       o.created_at AS order_created,
       jsonb_agg(
               jsonb_build_object(
                       'suborder_id', so.id,
                       'merchant_id', so.merchant_id,
                       'total_amount', so.total_amount,
                       'currency', so.currency,
                       'status', so.status,
                       'items', so.items,
                       'created_at', so.created_at,
                       'updated_at', so.updated_at
               ) ORDER BY so.created_at
       )            AS suborders
FROM orders.orders o
         LEFT JOIN orders.sub_orders so ON o.id = so.order_id
WHERE o.user_id = $1::uuid
GROUP BY o.id
ORDER BY o.created_at DESC
`

type GetUserOrdersWithSubordersRow struct {
	OrderID       uuid.UUID `json:"orderID"`
	OrderCurrency string    `json:"orderCurrency"`
	StreetAddress string    `json:"streetAddress"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	Country       string    `json:"country"`
	ZipCode       string    `json:"zipCode"`
	Email         string    `json:"email"`
	OrderCreated  time.Time `json:"orderCreated"`
	Suborders     []byte    `json:"suborders"`
}

// -- name: ListOrdersByUser :many
// SELECT *
// FROM orders.orders
// WHERE user_id = $1
// ORDER BY created_at DESC
// LIMIT $2 OFFSET $3;
//
//	SELECT o.id         AS order_id,
//	       o.currency   AS order_currency,
//	       o.street_address,
//	       o.city,
//	       o.state,
//	       o.country,
//	       o.zip_code,
//	       o.email,
//	       o.created_at AS order_created,
//	       jsonb_agg(
//	               jsonb_build_object(
//	                       'suborder_id', so.id,
//	                       'merchant_id', so.merchant_id,
//	                       'total_amount', so.total_amount,
//	                       'currency', so.currency,
//	                       'status', so.status,
//	                       'items', so.items,
//	                       'created_at', so.created_at,
//	                       'updated_at', so.updated_at
//	               ) ORDER BY so.created_at
//	       )            AS suborders
//	FROM orders.orders o
//	         LEFT JOIN orders.sub_orders so ON o.id = so.order_id
//	WHERE o.user_id = $1::uuid
//	GROUP BY o.id
//	ORDER BY o.created_at DESC
func (q *Queries) GetUserOrdersWithSuborders(ctx context.Context, dollar_1 uuid.UUID) ([]GetUserOrdersWithSubordersRow, error) {
	rows, err := q.db.Query(ctx, GetUserOrdersWithSuborders, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserOrdersWithSubordersRow
	for rows.Next() {
		var i GetUserOrdersWithSubordersRow
		if err := rows.Scan(
			&i.OrderID,
			&i.OrderCurrency,
			&i.StreetAddress,
			&i.City,
			&i.State,
			&i.Country,
			&i.ZipCode,
			&i.Email,
			&i.OrderCreated,
			&i.Suborders,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const UpdateSubOrderStatus = `-- name: UpdateSubOrderStatus :exec
UPDATE orders.sub_orders
SET status     = $2,
    updated_at = $3
WHERE id = $1
`

type UpdateSubOrderStatusParams struct {
	ID        uuid.UUID          `json:"id"`
	Status    string             `json:"status"`
	UpdatedAt pgtype.Timestamptz `json:"updatedAt"`
}

// UpdateSubOrderStatus
//
//	UPDATE orders.sub_orders
//	SET status     = $2,
//	    updated_at = $3
//	WHERE id = $1
func (q *Queries) UpdateSubOrderStatus(ctx context.Context, arg UpdateSubOrderStatusParams) error {
	_, err := q.db.Exec(ctx, UpdateSubOrderStatus, arg.ID, arg.Status, arg.UpdatedAt)
	return err
}
