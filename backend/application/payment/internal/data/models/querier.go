// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package models

import (
	"context"
)

type Querier interface {
	//CreatePayRecord
	//
	//  INSERT INTO pay_record (
	//      user_id,
	//      order_id,
	//      transcation_id,
	//      amount,
	//      pay_at,
	//      status
	//  ) VALUES (
	//      $1, $2, $3, $4, $5, $6
	//  ) RETURNING id, created_at, deleted_at, user_id, order_id, transcation_id, amount, pay_at, status
	CreatePayRecord(ctx context.Context, arg CreatePayRecordParams) (PayRecord, error)
	//GetPayRecordByOrderId
	//
	//  SELECT id, created_at, deleted_at, user_id, order_id, transcation_id, amount, pay_at, status
	//  FROM pay_record
	//  WHERE order_id = $1 AND deleted_at IS NULL
	//  LIMIT 1
	GetPayRecordByOrderId(ctx context.Context, orderID string) (PayRecord, error)
	//GetPayRecordByTransactionId
	//
	//  SELECT id, created_at, deleted_at, user_id, order_id, transcation_id, amount, pay_at, status
	//  FROM pay_record
	//  WHERE transcation_id = $1 AND deleted_at IS NULL
	//  LIMIT 1
	GetPayRecordByTransactionId(ctx context.Context, transcationID string) (PayRecord, error)
	//GetPayRecordsByUserId
	//
	//  SELECT id, created_at, deleted_at, user_id, order_id, transcation_id, amount, pay_at, status
	//  FROM pay_record
	//  WHERE user_id = $1 AND deleted_at IS NULL
	//  ORDER BY created_at DESC
	GetPayRecordsByUserId(ctx context.Context, userID string) ([]PayRecord, error)
}

var _ Querier = (*Queries)(nil)
