// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: product.sql

package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const BulkCreateProductImages = `-- name: BulkCreateProductImages :exec
INSERT INTO products.product_images
    (merchant_id, product_id, url, is_primary, sort_order)
SELECT m_id, p_id, u, is_p, s_ord
FROM ROWS FROM (
         unnest($1::bigint[]),
         unnest($2::bigint[]),
         unnest($3::text[]),
         unnest($4::boolean[]),
         unnest($5::smallint[])
         ) AS t(m_id, p_id, u, is_p, s_ord)
`

type BulkCreateProductImagesParams struct {
	MerchantIds []int64  `json:"merchantIds"`
	ProductIds  []int64  `json:"productIds"`
	Urls        []string `json:"urls"`
	IsPrimary   []bool   `json:"isPrimary"`
	SortOrders  []int16  `json:"sortOrders"`
}

// 批量插入图片
//
//	INSERT INTO products.product_images
//	    (merchant_id, product_id, url, is_primary, sort_order)
//	SELECT m_id, p_id, u, is_p, s_ord
//	FROM ROWS FROM (
//	         unnest($1::bigint[]),
//	         unnest($2::bigint[]),
//	         unnest($3::text[]),
//	         unnest($4::boolean[]),
//	         unnest($5::smallint[])
//	         ) AS t(m_id, p_id, u, is_p, s_ord)
func (q *Queries) BulkCreateProductImages(ctx context.Context, arg BulkCreateProductImagesParams) error {
	_, err := q.db.Exec(ctx, BulkCreateProductImages,
		arg.MerchantIds,
		arg.ProductIds,
		arg.Urls,
		arg.IsPrimary,
		arg.SortOrders,
	)
	return err
}

const CreateAuditRecord = `-- name: CreateAuditRecord :one
INSERT INTO product_audits (product_id,
                            merchant_id,
                            old_status,
                            new_status,
                            reason,
                            operator_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at
`

type CreateAuditRecordParams struct {
	Column1 *int64  `json:"column1"`
	Column2 *int64  `json:"column2"`
	Column3 *int16  `json:"column3"`
	Column4 *int16  `json:"column4"`
	Column5 *string `json:"column5"`
	Column6 *int64  `json:"column6"`
}

type CreateAuditRecordRow struct {
	ID        int64              `json:"id"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
}

// 创建审核记录，返回新记录ID
//
//	INSERT INTO product_audits (product_id,
//	                            merchant_id,
//	                            old_status,
//	                            new_status,
//	                            reason,
//	                            operator_id)
//	VALUES ($1, $2, $3, $4, $5, $6)
//	RETURNING id, created_at
func (q *Queries) CreateAuditRecord(ctx context.Context, arg CreateAuditRecordParams) (CreateAuditRecordRow, error) {
	row := q.db.QueryRow(ctx, CreateAuditRecord,
		arg.Column1,
		arg.Column2,
		arg.Column3,
		arg.Column4,
		arg.Column5,
		arg.Column6,
	)
	var i CreateAuditRecordRow
	err := row.Scan(&i.ID, &i.CreatedAt)
	return i, err
}

const CreateProduct = `-- name: CreateProduct :one

INSERT INTO products.products (name,
                               description,
                               price,
                               stock,
                               status,
                               merchant_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at, updated_at
`

type CreateProductParams struct {
	Name        string         `json:"name"`
	Description *string        `json:"description"`
	Price       pgtype.Numeric `json:"price"`
	Stock       *int32         `json:"stock"`
	Status      int16          `json:"status"`
	MerchantID  int64          `json:"merchantID"`
}

type CreateProductRow struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// 所有分片表必须：
// 1. 包含分片键列（merchant_id）
// 2. 主键必须包含分片键
// 3. 外键约束需要特殊处理（Citus 不支持跨节点外键）
// 创建商品主记录，返回生成的ID
// merchant_id 作为分片键，必须提供
//
//	INSERT INTO products.products (name,
//	                               description,
//	                               price,
//	                               stock,
//	                               status,
//	                               merchant_id)
//	VALUES ($1, $2, $3, $4, $5, $6)
//	RETURNING id, created_at, updated_at
func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (CreateProductRow, error) {
	row := q.db.QueryRow(ctx, CreateProduct,
		arg.Name,
		arg.Description,
		arg.Price,
		arg.Stock,
		arg.Status,
		arg.MerchantID,
	)
	var i CreateProductRow
	err := row.Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

type CreateProductImagesParams struct {
	MerchantID int64  `json:"merchantID"`
	ProductID  int64  `json:"productID"`
	Url        string `json:"url"`
	IsPrimary  bool   `json:"isPrimary"`
	SortOrder  *int16 `json:"sortOrder"`
}

const GetLatestAudit = `-- name: GetLatestAudit :one
INSERT INTO products.product_audits (merchant_id, -- 新增分片键
                                     product_id,
                                     old_status,
                                     new_status,
                                     reason,
                                     operator_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at
`

type GetLatestAuditParams struct {
	MerchantID int64   `json:"merchantID"`
	ProductID  int64   `json:"productID"`
	OldStatus  int16   `json:"oldStatus"`
	NewStatus  int16   `json:"newStatus"`
	Reason     *string `json:"reason"`
	OperatorID int64   `json:"operatorID"`
}

type GetLatestAuditRow struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

// 获取最新审核记录
//
//	INSERT INTO products.product_audits (merchant_id, -- 新增分片键
//	                                     product_id,
//	                                     old_status,
//	                                     new_status,
//	                                     reason,
//	                                     operator_id)
//	VALUES ($1, $2, $3, $4, $5, $6)
//	RETURNING id, created_at
func (q *Queries) GetLatestAudit(ctx context.Context, arg GetLatestAuditParams) (GetLatestAuditRow, error) {
	row := q.db.QueryRow(ctx, GetLatestAudit,
		arg.MerchantID,
		arg.ProductID,
		arg.OldStatus,
		arg.NewStatus,
		arg.Reason,
		arg.OperatorID,
	)
	var i GetLatestAuditRow
	err := row.Scan(&i.ID, &i.CreatedAt)
	return i, err
}

const GetProduct = `-- name: GetProduct :one

SELECT id,
       name,
       description,
       price,
       stock,
       status,
       merchant_id,
       created_at,
       updated_at
FROM products.products
WHERE id = $1
  AND merchant_id = $2
  AND deleted_at IS NULL
`

type GetProductParams struct {
	ID         int64 `json:"id"`
	MerchantID int64 `json:"merchantID"`
}

type GetProductRow struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Description *string        `json:"description"`
	Price       pgtype.Numeric `json:"price"`
	Stock       *int32         `json:"stock"`
	Status      int16          `json:"status"`
	MerchantID  int64          `json:"merchantID"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

// 乐观锁版本控制
// 获取商品详情，包含软删除检查
//
//	SELECT id,
//	       name,
//	       description,
//	       price,
//	       stock,
//	       status,
//	       merchant_id,
//	       created_at,
//	       updated_at
//	FROM products.products
//	WHERE id = $1
//	  AND merchant_id = $2
//	  AND deleted_at IS NULL
func (q *Queries) GetProduct(ctx context.Context, arg GetProductParams) (GetProductRow, error) {
	row := q.db.QueryRow(ctx, GetProduct, arg.ID, arg.MerchantID)
	var i GetProductRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.Stock,
		&i.Status,
		&i.MerchantID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const GetProductImages = `-- name: GetProductImages :many
SELECT id, merchant_id, product_id, url, is_primary, sort_order, created_at
FROM products.product_images
WHERE merchant_id = $1
  AND product_id = $2 -- 查询必须包含分片键
ORDER BY sort_order
`

type GetProductImagesParams struct {
	MerchantID int64 `json:"merchantID"`
	ProductID  int64 `json:"productID"`
}

// 获取商品图片列表，按排序顺序返回
//
//	SELECT id, merchant_id, product_id, url, is_primary, sort_order, created_at
//	FROM products.product_images
//	WHERE merchant_id = $1
//	  AND product_id = $2 -- 查询必须包含分片键
//	ORDER BY sort_order
func (q *Queries) GetProductImages(ctx context.Context, arg GetProductImagesParams) ([]ProductsProductImages, error) {
	rows, err := q.db.Query(ctx, GetProductImages, arg.MerchantID, arg.ProductID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ProductsProductImages
	for rows.Next() {
		var i ProductsProductImages
		if err := rows.Scan(
			&i.ID,
			&i.MerchantID,
			&i.ProductID,
			&i.Url,
			&i.IsPrimary,
			&i.SortOrder,
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

const SoftDeleteProduct = `-- name: SoftDeleteProduct :exec
UPDATE products.products
SET deleted_at = NOW()
WHERE id = $1
  AND merchant_id = $2
`

type SoftDeleteProductParams struct {
	ID         int64 `json:"id"`
	MerchantID int64 `json:"merchantID"`
}

// 软删除商品，设置删除时间戳
//
//	UPDATE products.products
//	SET deleted_at = NOW()
//	WHERE id = $1
//	  AND merchant_id = $2
func (q *Queries) SoftDeleteProduct(ctx context.Context, arg SoftDeleteProductParams) error {
	_, err := q.db.Exec(ctx, SoftDeleteProduct, arg.ID, arg.MerchantID)
	return err
}

const UpdateProduct = `-- name: UpdateProduct :exec
UPDATE products.products
SET name        = $2,
    description = $3,
    price       = $4,
    stock       = $5,
    status      = $6,
    updated_at  = NOW()
WHERE id = $1
  AND merchant_id = $7
  AND updated_at = $8
`

type UpdateProductParams struct {
	ID          int64              `json:"id"`
	Name        string             `json:"name"`
	Description *string            `json:"description"`
	Price       pgtype.Numeric     `json:"price"`
	Stock       *int32             `json:"stock"`
	Status      int16              `json:"status"`
	MerchantID  int64              `json:"merchantID"`
	UpdatedAt   pgtype.Timestamptz `json:"updatedAt"`
}

// 更新商品基础信息，使用乐观锁控制并发
//
//	UPDATE products.products
//	SET name        = $2,
//	    description = $3,
//	    price       = $4,
//	    stock       = $5,
//	    status      = $6,
//	    updated_at  = NOW()
//	WHERE id = $1
//	  AND merchant_id = $7
//	  AND updated_at = $8
func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) error {
	_, err := q.db.Exec(ctx, UpdateProduct,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.Price,
		arg.Stock,
		arg.Status,
		arg.MerchantID,
		arg.UpdatedAt,
	)
	return err
}

const UpdateProductStatus = `-- name: UpdateProductStatus :exec
UPDATE products.products
SET status           = $2,
    current_audit_id = $3,
    updated_at       = NOW()
WHERE id = $1
  AND merchant_id = $4
`

type UpdateProductStatusParams struct {
	ID             int64  `json:"id"`
	Status         int16  `json:"status"`
	CurrentAuditID *int64 `json:"currentAuditID"`
	MerchantID     int64  `json:"merchantID"`
}

// 更新商品状态并记录当前审核ID
//
//	UPDATE products.products
//	SET status           = $2,
//	    current_audit_id = $3,
//	    updated_at       = NOW()
//	WHERE id = $1
//	  AND merchant_id = $4
func (q *Queries) UpdateProductStatus(ctx context.Context, arg UpdateProductStatusParams) error {
	_, err := q.db.Exec(ctx, UpdateProductStatus,
		arg.ID,
		arg.Status,
		arg.CurrentAuditID,
		arg.MerchantID,
	)
	return err
}
