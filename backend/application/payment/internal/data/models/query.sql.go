// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package models

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const CreatePaymentQuery = `-- name: CreatePaymentQuery :one
INSERT INTO payments.payments (id, order_id, user_id, amount, currency, method, status,
                               subject, trade_no, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, order_id, user_id, amount, currency, method, status, subject, trade_no, metadata, created_at, updated_at
`

type CreatePaymentQueryParams struct {
	ID       int64          `json:"id"`
	OrderID  int64          `json:"orderID"`
	UserID   uuid.UUID      `json:"userID"`
	Amount   pgtype.Numeric `json:"amount"`
	Currency string         `json:"currency"`
	Method   string         `json:"method"`
	Status   string         `json:"status"`
	Subject  string         `json:"subject"`
	TradeNo  string         `json:"tradeNo"`
	Metadata []byte         `json:"metadata"`
}

// CreatePaymentQuery
//
//	INSERT INTO payments.payments (id, order_id, user_id, amount, currency, method, status,
//	                               subject, trade_no, metadata)
//	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
//	RETURNING id, order_id, user_id, amount, currency, method, status, subject, trade_no, metadata, created_at, updated_at
func (q *Queries) CreatePaymentQuery(ctx context.Context, arg CreatePaymentQueryParams) (PaymentsPayments, error) {
	row := q.db.QueryRow(ctx, CreatePaymentQuery,
		arg.ID,
		arg.OrderID,
		arg.UserID,
		arg.Amount,
		arg.Currency,
		arg.Method,
		arg.Status,
		arg.Subject,
		arg.TradeNo,
		arg.Metadata,
	)
	var i PaymentsPayments
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.UserID,
		&i.Amount,
		&i.Currency,
		&i.Method,
		&i.Status,
		&i.Subject,
		&i.TradeNo,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const GetByIDQuery = `-- name: GetByIDQuery :one
SELECT id, order_id, user_id, amount, currency, method, status, subject, trade_no, metadata, created_at, updated_at
FROM payments.payments
WHERE id = $1
`

// GetByIDQuery
//
//	SELECT id, order_id, user_id, amount, currency, method, status, subject, trade_no, metadata, created_at, updated_at
//	FROM payments.payments
//	WHERE id = $1
func (q *Queries) GetByIDQuery(ctx context.Context, id int64) (PaymentsPayments, error) {
	row := q.db.QueryRow(ctx, GetByIDQuery, id)
	var i PaymentsPayments
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.UserID,
		&i.Amount,
		&i.Currency,
		&i.Method,
		&i.Status,
		&i.Subject,
		&i.TradeNo,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const GetByOrderIDQuery = `-- name: GetByOrderIDQuery :one
SELECT id, order_id, user_id, amount, currency, method, status, subject, trade_no, metadata, created_at, updated_at
FROM payments.payments
WHERE order_id = $1
`

// GetByOrderIDQuery
//
//	SELECT id, order_id, user_id, amount, currency, method, status, subject, trade_no, metadata, created_at, updated_at
//	FROM payments.payments
//	WHERE order_id = $1
func (q *Queries) GetByOrderIDQuery(ctx context.Context, orderID int64) (PaymentsPayments, error) {
	row := q.db.QueryRow(ctx, GetByOrderIDQuery, orderID)
	var i PaymentsPayments
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.UserID,
		&i.Amount,
		&i.Currency,
		&i.Method,
		&i.Status,
		&i.Subject,
		&i.TradeNo,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const GetPayment = `-- name: GetPayment :one
SELECT id, order_id, user_id, amount, currency, method, status, subject, trade_no, metadata, created_at, updated_at
FROM payments.payments
WHERE id = $1
`

// GetPayment
//
//	SELECT id, order_id, user_id, amount, currency, method, status, subject, trade_no, metadata, created_at, updated_at
//	FROM payments.payments
//	WHERE id = $1
func (q *Queries) GetPayment(ctx context.Context, id int64) (PaymentsPayments, error) {
	row := q.db.QueryRow(ctx, GetPayment, id)
	var i PaymentsPayments
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.UserID,
		&i.Amount,
		&i.Currency,
		&i.Method,
		&i.Status,
		&i.Subject,
		&i.TradeNo,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const GetPaymentByOrderID = `-- name: GetPaymentByOrderID :one
SELECT id, order_id, user_id, amount, currency, method, status, subject, trade_no, metadata, created_at, updated_at
FROM payments.payments
WHERE order_id = $1
`

// GetPaymentByOrderID
//
//	SELECT id, order_id, user_id, amount, currency, method, status, subject, trade_no, metadata, created_at, updated_at
//	FROM payments.payments
//	WHERE order_id = $1
func (q *Queries) GetPaymentByOrderID(ctx context.Context, orderID int64) (PaymentsPayments, error) {
	row := q.db.QueryRow(ctx, GetPaymentByOrderID, orderID)
	var i PaymentsPayments
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.UserID,
		&i.Amount,
		&i.Currency,
		&i.Method,
		&i.Status,
		&i.Subject,
		&i.TradeNo,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const GetPaymentByTradeNo = `-- name: GetPaymentByTradeNo :one
SELECT id, order_id, user_id, amount, currency, method, status, subject, trade_no, metadata, created_at, updated_at
FROM payments.payments
WHERE trade_no = $1
`

// GetPaymentByTradeNo
//
//	SELECT id, order_id, user_id, amount, currency, method, status, subject, trade_no, metadata, created_at, updated_at
//	FROM payments.payments
//	WHERE trade_no = $1
func (q *Queries) GetPaymentByTradeNo(ctx context.Context, tradeNo string) (PaymentsPayments, error) {
	row := q.db.QueryRow(ctx, GetPaymentByTradeNo, tradeNo)
	var i PaymentsPayments
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.UserID,
		&i.Amount,
		&i.Currency,
		&i.Method,
		&i.Status,
		&i.Subject,
		&i.TradeNo,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const UpdatePaymentStatus = `-- name: UpdatePaymentStatus :one
UPDATE payments.payments
SET status     = $1,
    updated_at = now()
WHERE (id = $2 AND $2 != 0)
   OR (order_id = $3 AND $3 != 0)
RETURNING id, order_id, user_id, amount, currency, method, status, subject, trade_no, metadata, created_at, updated_at
`

type UpdatePaymentStatusParams struct {
	Status  string `json:"status"`
	ID      int64  `json:"id"`
	OrderID int64  `json:"orderID"`
}

// UpdatePaymentStatus
//
//	UPDATE payments.payments
//	SET status     = $1,
//	    updated_at = now()
//	WHERE (id = $2 AND $2 != 0)
//	   OR (order_id = $3 AND $3 != 0)
//	RETURNING id, order_id, user_id, amount, currency, method, status, subject, trade_no, metadata, created_at, updated_at
func (q *Queries) UpdatePaymentStatus(ctx context.Context, arg UpdatePaymentStatusParams) (PaymentsPayments, error) {
	row := q.db.QueryRow(ctx, UpdatePaymentStatus, arg.Status, arg.ID, arg.OrderID)
	var i PaymentsPayments
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.UserID,
		&i.Amount,
		&i.Currency,
		&i.Method,
		&i.Status,
		&i.Subject,
		&i.TradeNo,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const UpdateStatusQuery = `-- name: UpdateStatusQuery :one
UPDATE payments.payments
SET status     = $2,
    id         = $3,
    updated_at = now()
WHERE id = $1
RETURNING id, order_id, user_id, amount, currency, method, status, subject, trade_no, metadata, created_at, updated_at
`

type UpdateStatusQueryParams struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
	ID_2   int64  `json:"id2"`
}

// UpdateStatusQuery
//
//	UPDATE payments.payments
//	SET status     = $2,
//	    id         = $3,
//	    updated_at = now()
//	WHERE id = $1
//	RETURNING id, order_id, user_id, amount, currency, method, status, subject, trade_no, metadata, created_at, updated_at
func (q *Queries) UpdateStatusQuery(ctx context.Context, arg UpdateStatusQueryParams) (PaymentsPayments, error) {
	row := q.db.QueryRow(ctx, UpdateStatusQuery, arg.ID, arg.Status, arg.ID_2)
	var i PaymentsPayments
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.UserID,
		&i.Amount,
		&i.Currency,
		&i.Method,
		&i.Status,
		&i.Subject,
		&i.TradeNo,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
