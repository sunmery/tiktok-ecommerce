// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type PaymentsPayments struct {
	ID        int64          `json:"id"`
	OrderID   int64          `json:"orderID"`
	UserID    uuid.UUID      `json:"userID"`
	Amount    pgtype.Numeric `json:"amount"`
	Currency  string         `json:"currency"`
	Method    string         `json:"method"`
	Status    string         `json:"status"`
	Subject   string         `json:"subject"`
	TradeNo   string         `json:"tradeNo"`
	Metadata  []byte         `json:"metadata"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}
