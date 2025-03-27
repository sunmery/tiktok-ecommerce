// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: inventory.sql

package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const CountLowStockProducts = `-- name: CountLowStockProducts :one
SELECT COUNT(*)
FROM products.products p
         JOIN products.inventory i
              ON p.id = i.product_id
                  AND p.merchant_id = i.merchant_id
         LEFT JOIN merchant.stock_alerts sa
                   ON p.id = sa.product_id
                       AND p.merchant_id = sa.merchant_id
WHERE i.stock <= COALESCE(sa.threshold, $1)
  AND p.merchant_id = $2::uuid
`

type CountLowStockProductsParams struct {
	Threshold  *int32
	MerchantID pgtype.UUID
}

// 获取低库存产品总数
//
//	SELECT COUNT(*)
//	FROM products.products p
//	         JOIN products.inventory i
//	              ON p.id = i.product_id
//	                  AND p.merchant_id = i.merchant_id
//	         LEFT JOIN merchant.stock_alerts sa
//	                   ON p.id = sa.product_id
//	                       AND p.merchant_id = sa.merchant_id
//	WHERE i.stock <= COALESCE(sa.threshold, $1)
//	  AND p.merchant_id = $2::uuid
func (q *Queries) CountLowStockProducts(ctx context.Context, arg CountLowStockProductsParams) (*int64, error) {
	row := q.db.QueryRow(ctx, CountLowStockProducts, arg.Threshold, arg.MerchantID)
	var count *int64
	err := row.Scan(&count)
	return count, err
}

const CountStockAdjustmentHistory = `-- name: CountStockAdjustmentHistory :one
SELECT COUNT(*)
FROM merchant.stock_adjustments sa
WHERE sa.product_id = $1::uuid
  AND sa.merchant_id = $2::uuid
`

type CountStockAdjustmentHistoryParams struct {
	ProductID  uuid.UUID
	MerchantID uuid.UUID
}

// 获取库存调整历史总数
//
//	SELECT COUNT(*)
//	FROM merchant.stock_adjustments sa
//	WHERE sa.product_id = $1::uuid
//	  AND sa.merchant_id = $2::uuid
func (q *Queries) CountStockAdjustmentHistory(ctx context.Context, arg CountStockAdjustmentHistoryParams) (int64, error) {
	row := q.db.QueryRow(ctx, CountStockAdjustmentHistory, arg.ProductID, arg.MerchantID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const CountStockAlerts = `-- name: CountStockAlerts :one
SELECT COUNT(*)
FROM merchant.stock_alerts sa
WHERE sa.merchant_id = $1::uuid
`

// 获取库存警报配置总数
//
//	SELECT COUNT(*)
//	FROM merchant.stock_alerts sa
//	WHERE sa.merchant_id = $1::uuid
func (q *Queries) CountStockAlerts(ctx context.Context, merchantID uuid.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, CountStockAlerts, merchantID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const GetLowStockProducts = `-- name: GetLowStockProducts :many
SELECT p.id                       as product_id,
       p.merchant_id,
       p.name                     as product_name,
       i.stock                    as current_stock,
       COALESCE(sa.threshold, 10) as threshold,
       COALESCE(
               (SELECT pi.url
                FROM products.product_images pi
                WHERE pi.product_id = p.id
                  AND pi.merchant_id = p.merchant_id
                  AND pi.is_primary = true
                LIMIT 1),
               (SELECT pi.url
                FROM products.product_images pi
                WHERE pi.product_id = p.id
                  AND pi.merchant_id = p.merchant_id
                LIMIT 1)
       )                          as image_url
FROM products.products p
         JOIN products.inventory i
              ON p.id = i.product_id
                  AND p.merchant_id = i.merchant_id
         LEFT JOIN merchant.stock_alerts sa
                   ON p.id = sa.product_id
                       AND p.merchant_id = sa.merchant_id
WHERE i.stock <= COALESCE(sa.threshold, 10)
  AND p.merchant_id = $1::UUID
ORDER BY (i.stock * 1.0 / COALESCE(sa.threshold, 10))
LIMIT $3 OFFSET $2
`

type GetLowStockProductsParams struct {
	MerchantID pgtype.UUID
	Page       *int64
	PageSize   *int64
}

type GetLowStockProductsRow struct {
	ProductID    uuid.UUID
	MerchantID   uuid.UUID
	ProductName  string
	CurrentStock int32
	Threshold    *int32
	ImageUrl     *string
}

// 获取低库存产品列表
//
//	SELECT p.id                       as product_id,
//	       p.merchant_id,
//	       p.name                     as product_name,
//	       i.stock                    as current_stock,
//	       COALESCE(sa.threshold, 10) as threshold,
//	       COALESCE(
//	               (SELECT pi.url
//	                FROM products.product_images pi
//	                WHERE pi.product_id = p.id
//	                  AND pi.merchant_id = p.merchant_id
//	                  AND pi.is_primary = true
//	                LIMIT 1),
//	               (SELECT pi.url
//	                FROM products.product_images pi
//	                WHERE pi.product_id = p.id
//	                  AND pi.merchant_id = p.merchant_id
//	                LIMIT 1)
//	       )                          as image_url
//	FROM products.products p
//	         JOIN products.inventory i
//	              ON p.id = i.product_id
//	                  AND p.merchant_id = i.merchant_id
//	         LEFT JOIN merchant.stock_alerts sa
//	                   ON p.id = sa.product_id
//	                       AND p.merchant_id = sa.merchant_id
//	WHERE i.stock <= COALESCE(sa.threshold, 10)
//	  AND p.merchant_id = $1::UUID
//	ORDER BY (i.stock * 1.0 / COALESCE(sa.threshold, 10))
//	LIMIT $3 OFFSET $2
func (q *Queries) GetLowStockProducts(ctx context.Context, arg GetLowStockProductsParams) ([]GetLowStockProductsRow, error) {
	rows, err := q.db.Query(ctx, GetLowStockProducts, arg.MerchantID, arg.Page, arg.PageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetLowStockProductsRow
	for rows.Next() {
		var i GetLowStockProductsRow
		if err := rows.Scan(
			&i.ProductID,
			&i.MerchantID,
			&i.ProductName,
			&i.CurrentStock,
			&i.Threshold,
			&i.ImageUrl,
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

const GetProductStock = `-- name: GetProductStock :one
SELECT p.id                                                                     as product_id,
       p.merchant_id,
       i.stock,
       COALESCE(sa.threshold, 0)                                                as alert_threshold,
       CASE WHEN i.stock <= COALESCE(sa.threshold, 10) THEN true ELSE false END as is_low_stock
FROM products.products p
         JOIN products.inventory i
              ON p.id = i.product_id
                  AND p.merchant_id = i.merchant_id
         LEFT JOIN merchant.stock_alerts sa
                   ON p.id = sa.product_id
                       AND p.merchant_id = sa.merchant_id
WHERE p.id = $1::uuid
  AND p.merchant_id = $2::uuid
`

type GetProductStockParams struct {
	ProductID  pgtype.UUID
	MerchantID pgtype.UUID
}

type GetProductStockRow struct {
	ProductID      uuid.UUID
	MerchantID     uuid.UUID
	Stock          int32
	AlertThreshold *int32
	IsLowStock     *bool
}

// 获取产品库存
//
//	SELECT p.id                                                                     as product_id,
//	       p.merchant_id,
//	       i.stock,
//	       COALESCE(sa.threshold, 0)                                                as alert_threshold,
//	       CASE WHEN i.stock <= COALESCE(sa.threshold, 10) THEN true ELSE false END as is_low_stock
//	FROM products.products p
//	         JOIN products.inventory i
//	              ON p.id = i.product_id
//	                  AND p.merchant_id = i.merchant_id
//	         LEFT JOIN merchant.stock_alerts sa
//	                   ON p.id = sa.product_id
//	                       AND p.merchant_id = sa.merchant_id
//	WHERE p.id = $1::uuid
//	  AND p.merchant_id = $2::uuid
func (q *Queries) GetProductStock(ctx context.Context, arg GetProductStockParams) (GetProductStockRow, error) {
	row := q.db.QueryRow(ctx, GetProductStock, arg.ProductID, arg.MerchantID)
	var i GetProductStockRow
	err := row.Scan(
		&i.ProductID,
		&i.MerchantID,
		&i.Stock,
		&i.AlertThreshold,
		&i.IsLowStock,
	)
	return i, err
}

const GetStockAdjustmentHistory = `-- name: GetStockAdjustmentHistory :many
SELECT sa.id,
       sa.product_id,
       sa.merchant_id,
       p.name as product_name,
       sa.quantity,
       sa.reason,
       sa.operator_id,
       sa.created_at
FROM merchant.stock_adjustments sa
         JOIN products.products p
              ON sa.product_id = p.id
                  AND sa.merchant_id = p.merchant_id
WHERE sa.merchant_id = $1::uuid
ORDER BY sa.created_at DESC
LIMIT $3 OFFSET $2
`

type GetStockAdjustmentHistoryParams struct {
	MerchantID pgtype.UUID
	Page       *int64
	PageSize   *int64
}

type GetStockAdjustmentHistoryRow struct {
	ID          uuid.UUID
	ProductID   uuid.UUID
	MerchantID  uuid.UUID
	ProductName string
	Quantity    int32
	Reason      *string
	OperatorID  uuid.UUID
	CreatedAt   time.Time
}

// 获取库存调整历史
// WHERE sa.product_id = @product_id::uuid
//
//	SELECT sa.id,
//	       sa.product_id,
//	       sa.merchant_id,
//	       p.name as product_name,
//	       sa.quantity,
//	       sa.reason,
//	       sa.operator_id,
//	       sa.created_at
//	FROM merchant.stock_adjustments sa
//	         JOIN products.products p
//	              ON sa.product_id = p.id
//	                  AND sa.merchant_id = p.merchant_id
//	WHERE sa.merchant_id = $1::uuid
//	ORDER BY sa.created_at DESC
//	LIMIT $3 OFFSET $2
func (q *Queries) GetStockAdjustmentHistory(ctx context.Context, arg GetStockAdjustmentHistoryParams) ([]GetStockAdjustmentHistoryRow, error) {
	rows, err := q.db.Query(ctx, GetStockAdjustmentHistory, arg.MerchantID, arg.Page, arg.PageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetStockAdjustmentHistoryRow
	for rows.Next() {
		var i GetStockAdjustmentHistoryRow
		if err := rows.Scan(
			&i.ID,
			&i.ProductID,
			&i.MerchantID,
			&i.ProductName,
			&i.Quantity,
			&i.Reason,
			&i.OperatorID,
			&i.CreatedAt,
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

const GetStockAlerts = `-- name: GetStockAlerts :many
SELECT sa.id,
       sa.product_id,
       sa.merchant_id,
       p.name  as product_name,
       i.stock as current_stock,
       sa.threshold,
       sa.created_at,
       sa.updated_at
FROM merchant.stock_alerts sa
         JOIN products.products p ON sa.product_id = p.id::uuid -- 显式转换
         JOIN products.inventory i ON sa.product_id = i.product_id::uuid
WHERE sa.merchant_id = $1::uuid -- 强制类型
ORDER BY sa.updated_at DESC
LIMIT $3 OFFSET $2
`

type GetStockAlertsParams struct {
	MerchantID pgtype.UUID
	Page       *int64
	PageSize   *int64
}

type GetStockAlertsRow struct {
	ID           uuid.UUID
	ProductID    uuid.UUID
	MerchantID   uuid.UUID
	ProductName  string
	CurrentStock int32
	Threshold    int32
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// 获取库存警报配置
//
//	SELECT sa.id,
//	       sa.product_id,
//	       sa.merchant_id,
//	       p.name  as product_name,
//	       i.stock as current_stock,
//	       sa.threshold,
//	       sa.created_at,
//	       sa.updated_at
//	FROM merchant.stock_alerts sa
//	         JOIN products.products p ON sa.product_id = p.id::uuid -- 显式转换
//	         JOIN products.inventory i ON sa.product_id = i.product_id::uuid
//	WHERE sa.merchant_id = $1::uuid -- 强制类型
//	ORDER BY sa.updated_at DESC
//	LIMIT $3 OFFSET $2
func (q *Queries) GetStockAlerts(ctx context.Context, arg GetStockAlertsParams) ([]GetStockAlertsRow, error) {
	rows, err := q.db.Query(ctx, GetStockAlerts, arg.MerchantID, arg.Page, arg.PageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetStockAlertsRow
	for rows.Next() {
		var i GetStockAlertsRow
		if err := rows.Scan(
			&i.ID,
			&i.ProductID,
			&i.MerchantID,
			&i.ProductName,
			&i.CurrentStock,
			&i.Threshold,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const RecordStockAdjustment = `-- name: RecordStockAdjustment :one
INSERT INTO merchant.stock_adjustments (product_id, merchant_id, quantity, reason, operator_id)
VALUES ($1::uuid, $2::uuid, $3, $4, $5::uuid)
RETURNING id, product_id, merchant_id, quantity, reason, operator_id, created_at
`

type RecordStockAdjustmentParams struct {
	ProductID  uuid.UUID
	MerchantID uuid.UUID
	Quantity   int32
	Reason     *string
	OperatorID uuid.UUID
}

// 记录库存调整
//
//	INSERT INTO merchant.stock_adjustments (product_id, merchant_id, quantity, reason, operator_id)
//	VALUES ($1::uuid, $2::uuid, $3, $4, $5::uuid)
//	RETURNING id, product_id, merchant_id, quantity, reason, operator_id, created_at
func (q *Queries) RecordStockAdjustment(ctx context.Context, arg RecordStockAdjustmentParams) (MerchantStockAdjustments, error) {
	row := q.db.QueryRow(ctx, RecordStockAdjustment,
		arg.ProductID,
		arg.MerchantID,
		arg.Quantity,
		arg.Reason,
		arg.OperatorID,
	)
	var i MerchantStockAdjustments
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.MerchantID,
		&i.Quantity,
		&i.Reason,
		&i.OperatorID,
		&i.CreatedAt,
	)
	return i, err
}

const SetStockAlert = `-- name: SetStockAlert :one
INSERT INTO merchant.stock_alerts (product_id, merchant_id, threshold)
VALUES ($1::uuid, $2::uuid, $3)
ON CONFLICT (product_id, merchant_id) DO UPDATE
    SET threshold  = $3,
        updated_at = NOW()
RETURNING id, merchant.stock_alerts.product_id, merchant_id, threshold, created_at, updated_at
`

type SetStockAlertParams struct {
	ProductID  uuid.UUID
	MerchantID uuid.UUID
	Threshold  int32
}

// 设置库存警报阈值
//
//	INSERT INTO merchant.stock_alerts (product_id, merchant_id, threshold)
//	VALUES ($1::uuid, $2::uuid, $3)
//	ON CONFLICT (product_id, merchant_id) DO UPDATE
//	    SET threshold  = $3,
//	        updated_at = NOW()
//	RETURNING id, merchant.stock_alerts.product_id, merchant_id, threshold, created_at, updated_at
func (q *Queries) SetStockAlert(ctx context.Context, arg SetStockAlertParams) (MerchantStockAlerts, error) {
	row := q.db.QueryRow(ctx, SetStockAlert, arg.ProductID, arg.MerchantID, arg.Threshold)
	var i MerchantStockAlerts
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.MerchantID,
		&i.Threshold,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const UpdateProductStock = `-- name: UpdateProductStock :one
UPDATE products.inventory
SET stock = stock + $1
WHERE product_id = $2::uuid
  AND merchant_id = $3::uuid
RETURNING product_id, merchant_id, stock
`

type UpdateProductStockParams struct {
	Stock      *int32
	ProductID  pgtype.UUID
	MerchantID pgtype.UUID
}

type UpdateProductStockRow struct {
	ProductID  uuid.UUID
	MerchantID uuid.UUID
	Stock      int32
}

// 更新产品库存
//
//	UPDATE products.inventory
//	SET stock = stock + $1
//	WHERE product_id = $2::uuid
//	  AND merchant_id = $3::uuid
//	RETURNING product_id, merchant_id, stock
func (q *Queries) UpdateProductStock(ctx context.Context, arg UpdateProductStockParams) (UpdateProductStockRow, error) {
	row := q.db.QueryRow(ctx, UpdateProductStock, arg.Stock, arg.ProductID, arg.MerchantID)
	var i UpdateProductStockRow
	err := row.Scan(&i.ProductID, &i.MerchantID, &i.Stock)
	return i, err
}
