// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// 商品库存表
type ProductsInventory struct {
	ProductID  uuid.UUID          `json:"productID"`
	MerchantID uuid.UUID          `json:"merchantID"`
	Stock      int32              `json:"stock"`
	CreatedAt  pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt  pgtype.Timestamptz `json:"updatedAt"`
}

// 商品属性表
type ProductsProductAttributes struct {
	MerchantID uuid.UUID `json:"merchantID"`
	ProductID  uuid.UUID `json:"productID"`
	Attributes []byte    `json:"attributes"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// 商品图片表
type ProductsProductImages struct {
	ID         uuid.UUID `json:"id"`
	MerchantID uuid.UUID `json:"merchantID"`
	ProductID  uuid.UUID `json:"productID"`
	Url        string    `json:"url"`
	IsPrimary  bool      `json:"isPrimary"`
	SortOrder  *int16    `json:"sortOrder"`
	CreatedAt  time.Time `json:"createdAt"`
}

// 商品表
type ProductsProducts struct {
	ID             uuid.UUID          `json:"id"`
	MerchantID     uuid.UUID          `json:"merchantID"`
	Name           string             `json:"name"`
	Description    *string            `json:"description"`
	Price          pgtype.Numeric     `json:"price"`
	Status         int16              `json:"status"`
	CurrentAuditID pgtype.UUID        `json:"currentAuditID"`
	CategoryID     int64              `json:"categoryID"`
	CreatedAt      time.Time          `json:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt"`
	DeletedAt      pgtype.Timestamptz `json:"deletedAt"`
}
