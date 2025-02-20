// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const CreatePayRecord = `-- name: CreatePayRecord :one
INSERT INTO pay_record (
    user_id,
    order_id,
    transcation_id,
    amount,
    pay_at,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING id, created_at, deleted_at, user_id, order_id, transcation_id, amount, pay_at, status
`

type CreatePayRecordParams struct {
	UserID        string             `json:"userID"`
	OrderID       string             `json:"orderID"`
	TranscationID string             `json:"transcationID"`
	Amount        float64            `json:"amount"`
	PayAt         pgtype.Timestamptz `json:"payAt"`
	Status        string             `json:"status"`
}

// CreatePayRecord
//
//	INSERT INTO pay_record (
//	    user_id,
//	    order_id,
//	    transcation_id,
//	    amount,
//	    pay_at,
//	    status
//	) VALUES (
//	    $1, $2, $3, $4, $5, $6
//	) RETURNING id, created_at, deleted_at, user_id, order_id, transcation_id, amount, pay_at, status
func (q *Queries) CreatePayRecord(ctx context.Context, arg CreatePayRecordParams) (PayRecord, error) {
	row := q.db.QueryRow(ctx, CreatePayRecord,
		arg.UserID,
		arg.OrderID,
		arg.TranscationID,
		arg.Amount,
		arg.PayAt,
		arg.Status,
	)
	var i PayRecord
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.DeletedAt,
		&i.UserID,
		&i.OrderID,
		&i.TranscationID,
		&i.Amount,
		&i.PayAt,
		&i.Status,
	)
	return i, err
}

const GetPayRecordByOrderId = `-- name: GetPayRecordByOrderId :one
SELECT id, created_at, deleted_at, user_id, order_id, transcation_id, amount, pay_at, status
FROM pay_record
WHERE order_id = $1 AND deleted_at IS NULL
LIMIT 1
`

// GetPayRecordByOrderId
//
//	SELECT id, created_at, deleted_at, user_id, order_id, transcation_id, amount, pay_at, status
//	FROM pay_record
//	WHERE order_id = $1 AND deleted_at IS NULL
//	LIMIT 1
func (q *Queries) GetPayRecordByOrderId(ctx context.Context, orderID string) (PayRecord, error) {
	row := q.db.QueryRow(ctx, GetPayRecordByOrderId, orderID)
	var i PayRecord
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.DeletedAt,
		&i.UserID,
		&i.OrderID,
		&i.TranscationID,
		&i.Amount,
		&i.PayAt,
		&i.Status,
	)
	return i, err
}

const GetPayRecordByTransactionId = `-- name: GetPayRecordByTransactionId :one
SELECT id, created_at, deleted_at, user_id, order_id, transcation_id, amount, pay_at, status
FROM pay_record
WHERE transcation_id = $1 AND deleted_at IS NULL
LIMIT 1
`

// GetPayRecordByTransactionId
//
//	SELECT id, created_at, deleted_at, user_id, order_id, transcation_id, amount, pay_at, status
//	FROM pay_record
//	WHERE transcation_id = $1 AND deleted_at IS NULL
//	LIMIT 1
func (q *Queries) GetPayRecordByTransactionId(ctx context.Context, transcationID string) (PayRecord, error) {
	row := q.db.QueryRow(ctx, GetPayRecordByTransactionId, transcationID)
	var i PayRecord
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.DeletedAt,
		&i.UserID,
		&i.OrderID,
		&i.TranscationID,
		&i.Amount,
		&i.PayAt,
		&i.Status,
	)
	return i, err
}

const GetPayRecordsByUserId = `-- name: GetPayRecordsByUserId :many
SELECT id, created_at, deleted_at, user_id, order_id, transcation_id, amount, pay_at, status
FROM pay_record
WHERE user_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
`

// GetPayRecordsByUserId
//
//	SELECT id, created_at, deleted_at, user_id, order_id, transcation_id, amount, pay_at, status
//	FROM pay_record
//	WHERE user_id = $1 AND deleted_at IS NULL
//	ORDER BY created_at DESC
func (q *Queries) GetPayRecordsByUserId(ctx context.Context, userID string) ([]PayRecord, error) {
	rows, err := q.db.Query(ctx, GetPayRecordsByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PayRecord
	for rows.Next() {
		var i PayRecord
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.DeletedAt,
			&i.UserID,
			&i.OrderID,
			&i.TranscationID,
			&i.Amount,
			&i.PayAt,
			&i.Status,
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
