// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type MerchantStockAdjustments struct {
	ID         uuid.UUID
	ProductID  uuid.UUID
	MerchantID uuid.UUID
	Quantity   uint32
	Reason     *string
	OperatorID uuid.UUID
	CreatedAt  pgtype.Timestamptz
}

type MerchantStockAlerts struct {
	ID         uuid.UUID
	ProductID  uuid.UUID
	MerchantID uuid.UUID
	Threshold  uint32
	CreatedAt  pgtype.Timestamptz
	UpdatedAt  pgtype.Timestamptz
}
