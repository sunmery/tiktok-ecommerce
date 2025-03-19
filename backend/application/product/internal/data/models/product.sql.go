// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: product.sql

package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const BulkCreateProductImages = `-- name: BulkCreateProductImages :exec
INSERT INTO products.product_images
    (merchant_id, product_id, url, is_primary, sort_order)
SELECT m_id, p_id, u, is_p, s_ord
FROM unnest(
             $1::uuid[],
             $2::uuid[],
             $3::text[],
             $4::boolean[],
             $5::smallint[]
     ) AS t(m_id, p_id, u, is_p, s_ord)
`

type BulkCreateProductImagesParams struct {
	MerchantIds []uuid.UUID `json:"merchantIds"`
	ProductIds  []uuid.UUID `json:"productIds"`
	Urls        []string    `json:"urls"`
	IsPrimary   []bool      `json:"isPrimary"`
	SortOrders  []int16     `json:"sortOrders"`
}

// 批量插入图片
//
//	INSERT INTO products.product_images
//	    (merchant_id, product_id, url, is_primary, sort_order)
//	SELECT m_id, p_id, u, is_p, s_ord
//	FROM unnest(
//	             $1::uuid[],
//	             $2::uuid[],
//	             $3::text[],
//	             $4::boolean[],
//	             $5::smallint[]
//	     ) AS t(m_id, p_id, u, is_p, s_ord)
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
INSERT INTO products.product_audits (product_id,
                                     merchant_id,
                                     old_status,
                                     new_status,
                                     reason,
                                     operator_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at
`

type CreateAuditRecordParams struct {
	ProductID  uuid.UUID `json:"productID"`
	MerchantID uuid.UUID `json:"merchantID"`
	OldStatus  int16     `json:"oldStatus"`
	NewStatus  int16     `json:"newStatus"`
	Reason     *string   `json:"reason"`
	OperatorID uuid.UUID `json:"operatorID"`
}

type CreateAuditRecordRow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

// 创建审核记录，返回新记录ID
//
//	INSERT INTO products.product_audits (product_id,
//	                                     merchant_id,
//	                                     old_status,
//	                                     new_status,
//	                                     reason,
//	                                     operator_id)
//	VALUES ($1, $2, $3, $4, $5, $6)
//	RETURNING id, created_at
func (q *Queries) CreateAuditRecord(ctx context.Context, arg CreateAuditRecordParams) (CreateAuditRecordRow, error) {
	row := q.db.QueryRow(ctx, CreateAuditRecord,
		arg.ProductID,
		arg.MerchantID,
		arg.OldStatus,
		arg.NewStatus,
		arg.Reason,
		arg.OperatorID,
	)
	var i CreateAuditRecordRow
	err := row.Scan(&i.ID, &i.CreatedAt)
	return i, err
}

const CreateInventory = `-- name: CreateInventory :one
INSERT INTO products.inventory (product_id, merchant_id, stock)
VALUES ($1, $2, $3)
RETURNING product_id, merchant_id, stock
`

type CreateInventoryParams struct {
	ProductID  uuid.UUID `json:"productID"`
	MerchantID uuid.UUID `json:"merchantID"`
	Stock      int32     `json:"stock"`
}

// CreateInventory
//
//	INSERT INTO products.inventory (product_id, merchant_id, stock)
//	VALUES ($1, $2, $3)
//	RETURNING product_id, merchant_id, stock
func (q *Queries) CreateInventory(ctx context.Context, arg CreateInventoryParams) (ProductsInventory, error) {
	row := q.db.QueryRow(ctx, CreateInventory, arg.ProductID, arg.MerchantID, arg.Stock)
	var i ProductsInventory
	err := row.Scan(&i.ProductID, &i.MerchantID, &i.Stock)
	return i, err
}

const CreateProduct = `-- name: CreateProduct :one

INSERT INTO products.products (name,
                               description,
                               price,
                               status,
                               merchant_id,
                               category_id
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at, updated_at
`

type CreateProductParams struct {
	Name        string         `json:"name"`
	Description *string        `json:"description"`
	Price       pgtype.Numeric `json:"price"`
	Status      int16          `json:"status"`
	MerchantID  uuid.UUID      `json:"merchantID"`
	CategoryID  int64          `json:"categoryID"`
}

type CreateProductRow struct {
	ID        uuid.UUID `json:"id"`
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
//	                               status,
//	                               merchant_id,
//	                               category_id
//	)
//	VALUES ($1, $2, $3, $4, $5, $6)
//	RETURNING id, created_at, updated_at
func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (CreateProductRow, error) {
	row := q.db.QueryRow(ctx, CreateProduct,
		arg.Name,
		arg.Description,
		arg.Price,
		arg.Status,
		arg.MerchantID,
		arg.CategoryID,
	)
	var i CreateProductRow
	err := row.Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}

const CreateProductAttribute = `-- name: CreateProductAttribute :exec
INSERT INTO products.product_attributes
    (merchant_id, product_id, attributes)
VALUES ($1, $2, $3)
RETURNING created_at, updated_at
`

type CreateProductAttributeParams struct {
	MerchantID uuid.UUID `json:"merchantID"`
	ProductID  uuid.UUID `json:"productID"`
	Attributes []byte    `json:"attributes"`
}

// CreateProductAttribute
//
//	INSERT INTO products.product_attributes
//	    (merchant_id, product_id, attributes)
//	VALUES ($1, $2, $3)
//	RETURNING created_at, updated_at
func (q *Queries) CreateProductAttribute(ctx context.Context, arg CreateProductAttributeParams) error {
	_, err := q.db.Exec(ctx, CreateProductAttribute, arg.MerchantID, arg.ProductID, arg.Attributes)
	return err
}

type CreateProductImagesParams struct {
	MerchantID uuid.UUID `json:"merchantID"`
	ProductID  uuid.UUID `json:"productID"`
	Url        string    `json:"url"`
	IsPrimary  bool      `json:"isPrimary"`
	SortOrder  *int16    `json:"sortOrder"`
}

const DeleteProductAttribute = `-- name: DeleteProductAttribute :exec
DELETE
FROM products.product_attributes
WHERE merchant_id = $1
  AND product_id = $2
`

type DeleteProductAttributeParams struct {
	MerchantID uuid.UUID `json:"merchantID"`
	ProductID  uuid.UUID `json:"productID"`
}

// DeleteProductAttribute
//
//	DELETE
//	FROM products.product_attributes
//	WHERE merchant_id = $1
//	  AND product_id = $2
func (q *Queries) DeleteProductAttribute(ctx context.Context, arg DeleteProductAttributeParams) error {
	_, err := q.db.Exec(ctx, DeleteProductAttribute, arg.MerchantID, arg.ProductID)
	return err
}

const GetCategoryProducts = `-- name: GetCategoryProducts :many
WITH filtered_products AS (SELECT p.id,
                                  p.merchant_id,
                                  p.name,
                                  p.description,
                                  p.price,
                                  p.status,
                                  p.category_id,
                                  p.created_at,
                                  p.updated_at
                           FROM products.products p
                           WHERE p.category_id = $1 -- 指定分类id
                             AND p.status = $2      -- 商品状态机
                             AND p.deleted_at IS NULL),
     product_images_agg AS (SELECT pi.product_id,
                                   jsonb_agg(
                                           jsonb_build_object(
                                                   'id', pi.id,
                                                   'url', pi.url,
                                                   'is_primary', pi.is_primary,
                                                   'sort_order', pi.sort_order
                                           )
                                   ) AS images
                            FROM products.product_images pi
                                     INNER JOIN filtered_products fp
                                                ON pi.product_id = fp.id AND pi.merchant_id = fp.merchant_id
                            GROUP BY pi.product_id),
     product_attributes_agg AS (SELECT pa.product_id,
                                       pa.attributes
                                FROM products.product_attributes pa
                                         INNER JOIN filtered_products fp
                                                    ON pa.product_id = fp.id AND pa.merchant_id = fp.merchant_id)
SELECT fp.id,
       fp.merchant_id,
       fp.name,
       fp.description,
       fp.price,
       fp.status,
       fp.category_id,
       fp.created_at,
       fp.updated_at,
       COALESCE(pia.images, '[]'::jsonb)     AS images,
       COALESCE(paa.attributes, '{}'::jsonb) AS attributes
FROM filtered_products fp
         LEFT JOIN product_images_agg pia
                   ON fp.id = pia.product_id
         LEFT JOIN product_attributes_agg paa
                   ON fp.id = paa.product_id
ORDER BY fp.created_at DESC
LIMIT $3 OFFSET $4
`

type GetCategoryProductsParams struct {
	CategoryID int64 `json:"categoryID"`
	Status     int16 `json:"status"`
	Limit      int64 `json:"limit"`
	Offset     int64 `json:"offset"`
}

type GetCategoryProductsRow struct {
	ID          uuid.UUID      `json:"id"`
	MerchantID  uuid.UUID      `json:"merchantID"`
	Name        string         `json:"name"`
	Description *string        `json:"description"`
	Price       pgtype.Numeric `json:"price"`
	Status      int16          `json:"status"`
	CategoryID  int64          `json:"categoryID"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	Images      []byte         `json:"images"`
	Attributes  []byte         `json:"attributes"`
}

// 根据分类获取商品列表
//
//	WITH filtered_products AS (SELECT p.id,
//	                                  p.merchant_id,
//	                                  p.name,
//	                                  p.description,
//	                                  p.price,
//	                                  p.status,
//	                                  p.category_id,
//	                                  p.created_at,
//	                                  p.updated_at
//	                           FROM products.products p
//	                           WHERE p.category_id = $1 -- 指定分类id
//	                             AND p.status = $2      -- 商品状态机
//	                             AND p.deleted_at IS NULL),
//	     product_images_agg AS (SELECT pi.product_id,
//	                                   jsonb_agg(
//	                                           jsonb_build_object(
//	                                                   'id', pi.id,
//	                                                   'url', pi.url,
//	                                                   'is_primary', pi.is_primary,
//	                                                   'sort_order', pi.sort_order
//	                                           )
//	                                   ) AS images
//	                            FROM products.product_images pi
//	                                     INNER JOIN filtered_products fp
//	                                                ON pi.product_id = fp.id AND pi.merchant_id = fp.merchant_id
//	                            GROUP BY pi.product_id),
//	     product_attributes_agg AS (SELECT pa.product_id,
//	                                       pa.attributes
//	                                FROM products.product_attributes pa
//	                                         INNER JOIN filtered_products fp
//	                                                    ON pa.product_id = fp.id AND pa.merchant_id = fp.merchant_id)
//	SELECT fp.id,
//	       fp.merchant_id,
//	       fp.name,
//	       fp.description,
//	       fp.price,
//	       fp.status,
//	       fp.category_id,
//	       fp.created_at,
//	       fp.updated_at,
//	       COALESCE(pia.images, '[]'::jsonb)     AS images,
//	       COALESCE(paa.attributes, '{}'::jsonb) AS attributes
//	FROM filtered_products fp
//	         LEFT JOIN product_images_agg pia
//	                   ON fp.id = pia.product_id
//	         LEFT JOIN product_attributes_agg paa
//	                   ON fp.id = paa.product_id
//	ORDER BY fp.created_at DESC
//	LIMIT $3 OFFSET $4
func (q *Queries) GetCategoryProducts(ctx context.Context, arg GetCategoryProductsParams) ([]GetCategoryProductsRow, error) {
	rows, err := q.db.Query(ctx, GetCategoryProducts,
		arg.CategoryID,
		arg.Status,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCategoryProductsRow
	for rows.Next() {
		var i GetCategoryProductsRow
		if err := rows.Scan(
			&i.ID,
			&i.MerchantID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.Status,
			&i.CategoryID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Images,
			&i.Attributes,
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
	MerchantID uuid.UUID `json:"merchantID"`
	ProductID  uuid.UUID `json:"productID"`
	OldStatus  int16     `json:"oldStatus"`
	NewStatus  int16     `json:"newStatus"`
	Reason     *string   `json:"reason"`
	OperatorID uuid.UUID `json:"operatorID"`
}

type GetLatestAuditRow struct {
	ID        uuid.UUID `json:"id"`
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

const GetMerchantProducts = `-- name: GetMerchantProducts :many
SELECT p.id,
       p.name,
       p.description,
       p.price,
       p.status,
       p.merchant_id,
       p.category_id,
       p.created_at,
       p.updated_at,
       i.stock,
       (SELECT jsonb_agg(jsonb_build_object(
               'url', pi.url,
               'is_primary', pi.is_primary,
               'sort_order', pi.sort_order
                         ))
        FROM products.product_images pi
        WHERE pi.merchant_id = p.merchant_id) AS images,
       pa.attributes,
       (SELECT jsonb_build_object(
                       'id', a.id,
                       'old_status', a.old_status,
                       'new_status', a.new_status,
                       'reason', a.reason,
                       'created_at', a.created_at
               )
        FROM products.product_audits a
        WHERE a.merchant_id = p.merchant_id
        ORDER BY a.created_at DESC
        LIMIT 1)                              AS latest_audit
FROM products.products p
         INNER JOIN products.inventory i
                    ON p.id = i.product_id AND p.merchant_id = i.merchant_id
         LEFT JOIN products.product_attributes pa
                   ON p.id = pa.product_id AND p.merchant_id = pa.merchant_id
WHERE p.merchant_id = $1
  AND p.deleted_at IS NULL
`

type GetMerchantProductsRow struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Description *string        `json:"description"`
	Price       pgtype.Numeric `json:"price"`
	Status      int16          `json:"status"`
	MerchantID  uuid.UUID      `json:"merchantID"`
	CategoryID  int64          `json:"categoryID"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	Stock       int32          `json:"stock"`
	Images      []byte         `json:"images"`
	Attributes  []byte         `json:"attributes"`
	LatestAudit []byte         `json:"latestAudit"`
}

// GetMerchantProducts
//
//	SELECT p.id,
//	       p.name,
//	       p.description,
//	       p.price,
//	       p.status,
//	       p.merchant_id,
//	       p.category_id,
//	       p.created_at,
//	       p.updated_at,
//	       i.stock,
//	       (SELECT jsonb_agg(jsonb_build_object(
//	               'url', pi.url,
//	               'is_primary', pi.is_primary,
//	               'sort_order', pi.sort_order
//	                         ))
//	        FROM products.product_images pi
//	        WHERE pi.merchant_id = p.merchant_id) AS images,
//	       pa.attributes,
//	       (SELECT jsonb_build_object(
//	                       'id', a.id,
//	                       'old_status', a.old_status,
//	                       'new_status', a.new_status,
//	                       'reason', a.reason,
//	                       'created_at', a.created_at
//	               )
//	        FROM products.product_audits a
//	        WHERE a.merchant_id = p.merchant_id
//	        ORDER BY a.created_at DESC
//	        LIMIT 1)                              AS latest_audit
//	FROM products.products p
//	         INNER JOIN products.inventory i
//	                    ON p.id = i.product_id AND p.merchant_id = i.merchant_id
//	         LEFT JOIN products.product_attributes pa
//	                   ON p.id = pa.product_id AND p.merchant_id = pa.merchant_id
//	WHERE p.merchant_id = $1
//	  AND p.deleted_at IS NULL
func (q *Queries) GetMerchantProducts(ctx context.Context, merchantID uuid.UUID) ([]GetMerchantProductsRow, error) {
	rows, err := q.db.Query(ctx, GetMerchantProducts, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMerchantProductsRow
	for rows.Next() {
		var i GetMerchantProductsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.Status,
			&i.MerchantID,
			&i.CategoryID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Stock,
			&i.Images,
			&i.Attributes,
			&i.LatestAudit,
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

const GetProduct = `-- name: GetProduct :one

SELECT p.id,
       p.name,
       p.description,
       p.price,
       p.status,
       p.merchant_id,
       p.category_id,
       p.created_at,
       p.updated_at,
       i.stock,
       (SELECT jsonb_agg(jsonb_build_object(
               'url', pi.url,
               'is_primary', pi.is_primary,
               'sort_order', pi.sort_order
                         ))
        FROM products.product_images pi
        WHERE pi.product_id = p.id
          AND pi.merchant_id = p.merchant_id) AS images,
       pa.attributes,
       (SELECT jsonb_build_object(
                       'id', a.id,
                       'old_status', a.old_status,
                       'new_status', a.new_status,
                       'reason', a.reason,
                       'created_at', a.created_at
               )
        FROM products.product_audits a
        WHERE a.product_id = p.id
          AND a.merchant_id = p.merchant_id
        ORDER BY a.created_at DESC
        LIMIT 1)                              AS latest_audit
FROM products.products p
         INNER JOIN products.inventory i
                    ON p.id = i.product_id AND p.merchant_id = i.merchant_id
         LEFT JOIN products.product_attributes pa
                   ON p.id = pa.product_id AND p.merchant_id = pa.merchant_id
WHERE p.id = $1
  AND p.merchant_id = $2
  AND p.deleted_at IS NULL
`

type GetProductParams struct {
	ID         uuid.UUID `json:"id"`
	MerchantID uuid.UUID `json:"merchantID"`
}

type GetProductRow struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Description *string        `json:"description"`
	Price       pgtype.Numeric `json:"price"`
	Status      int16          `json:"status"`
	MerchantID  uuid.UUID      `json:"merchantID"`
	CategoryID  int64          `json:"categoryID"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	Stock       int32          `json:"stock"`
	Images      []byte         `json:"images"`
	Attributes  []byte         `json:"attributes"`
	LatestAudit []byte         `json:"latestAudit"`
}

// 乐观锁版本控制
// 获取商品详情，包含软删除检查
//
//	SELECT p.id,
//	       p.name,
//	       p.description,
//	       p.price,
//	       p.status,
//	       p.merchant_id,
//	       p.category_id,
//	       p.created_at,
//	       p.updated_at,
//	       i.stock,
//	       (SELECT jsonb_agg(jsonb_build_object(
//	               'url', pi.url,
//	               'is_primary', pi.is_primary,
//	               'sort_order', pi.sort_order
//	                         ))
//	        FROM products.product_images pi
//	        WHERE pi.product_id = p.id
//	          AND pi.merchant_id = p.merchant_id) AS images,
//	       pa.attributes,
//	       (SELECT jsonb_build_object(
//	                       'id', a.id,
//	                       'old_status', a.old_status,
//	                       'new_status', a.new_status,
//	                       'reason', a.reason,
//	                       'created_at', a.created_at
//	               )
//	        FROM products.product_audits a
//	        WHERE a.product_id = p.id
//	          AND a.merchant_id = p.merchant_id
//	        ORDER BY a.created_at DESC
//	        LIMIT 1)                              AS latest_audit
//	FROM products.products p
//	         INNER JOIN products.inventory i
//	                    ON p.id = i.product_id AND p.merchant_id = i.merchant_id
//	         LEFT JOIN products.product_attributes pa
//	                   ON p.id = pa.product_id AND p.merchant_id = pa.merchant_id
//	WHERE p.id = $1
//	  AND p.merchant_id = $2
//	  AND p.deleted_at IS NULL
func (q *Queries) GetProduct(ctx context.Context, arg GetProductParams) (GetProductRow, error) {
	row := q.db.QueryRow(ctx, GetProduct, arg.ID, arg.MerchantID)
	var i GetProductRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.Status,
		&i.MerchantID,
		&i.CategoryID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Stock,
		&i.Images,
		&i.Attributes,
		&i.LatestAudit,
	)
	return i, err
}

const GetProductAttribute = `-- name: GetProductAttribute :one
SELECT merchant_id,
       product_id,
       attributes,
       created_at,
       updated_at
FROM products.product_attributes
WHERE merchant_id = $1
  AND product_id = $2
`

type GetProductAttributeParams struct {
	MerchantID uuid.UUID `json:"merchantID"`
	ProductID  uuid.UUID `json:"productID"`
}

// GetProductAttribute
//
//	SELECT merchant_id,
//	       product_id,
//	       attributes,
//	       created_at,
//	       updated_at
//	FROM products.product_attributes
//	WHERE merchant_id = $1
//	  AND product_id = $2
func (q *Queries) GetProductAttribute(ctx context.Context, arg GetProductAttributeParams) (ProductsProductAttributes, error) {
	row := q.db.QueryRow(ctx, GetProductAttribute, arg.MerchantID, arg.ProductID)
	var i ProductsProductAttributes
	err := row.Scan(
		&i.MerchantID,
		&i.ProductID,
		&i.Attributes,
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
	MerchantID uuid.UUID `json:"merchantID"`
	ProductID  uuid.UUID `json:"productID"`
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

const ListProductsByCategory = `-- name: ListProductsByCategory :many
SELECT p.id,
       p.merchant_id,
       p.name,
       p.description,
       p.price,
       p.status,
       p.category_id,
       p.created_at,
       p.updated_at,
       (SELECT jsonb_agg(jsonb_build_object(
               'url', pi.url,
               'is_primary', pi.is_primary,
               'sort_order', pi.sort_order
                         ))
        FROM products.product_images pi
        WHERE pi.product_id = p.id
          AND pi.merchant_id = p.merchant_id) AS images,
       pa.attributes
FROM products.products p
         LEFT JOIN products.product_attributes pa
                   ON p.id = pa.product_id AND p.merchant_id = pa.merchant_id
WHERE p.category_id = ANY ($1::bigint[])
  AND p.status = $2
  AND p.deleted_at IS NULL
ORDER BY p.created_at DESC
LIMIT $3 OFFSET $4
`

type ListProductsByCategoryParams struct {
	Column1 []int64 `json:"column1"`
	Status  int16   `json:"status"`
	Limit   int64   `json:"limit"`
	Offset  int64   `json:"offset"`
}

type ListProductsByCategoryRow struct {
	ID          uuid.UUID      `json:"id"`
	MerchantID  uuid.UUID      `json:"merchantID"`
	Name        string         `json:"name"`
	Description *string        `json:"description"`
	Price       pgtype.Numeric `json:"price"`
	Status      int16          `json:"status"`
	CategoryID  int64          `json:"categoryID"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	Images      []byte         `json:"images"`
	Attributes  []byte         `json:"attributes"`
}

// 分类批量查询（使用GIN索引优化数组查询）
//
//	SELECT p.id,
//	       p.merchant_id,
//	       p.name,
//	       p.description,
//	       p.price,
//	       p.status,
//	       p.category_id,
//	       p.created_at,
//	       p.updated_at,
//	       (SELECT jsonb_agg(jsonb_build_object(
//	               'url', pi.url,
//	               'is_primary', pi.is_primary,
//	               'sort_order', pi.sort_order
//	                         ))
//	        FROM products.product_images pi
//	        WHERE pi.product_id = p.id
//	          AND pi.merchant_id = p.merchant_id) AS images,
//	       pa.attributes
//	FROM products.products p
//	         LEFT JOIN products.product_attributes pa
//	                   ON p.id = pa.product_id AND p.merchant_id = pa.merchant_id
//	WHERE p.category_id = ANY ($1::bigint[])
//	  AND p.status = $2
//	  AND p.deleted_at IS NULL
//	ORDER BY p.created_at DESC
//	LIMIT $3 OFFSET $4
func (q *Queries) ListProductsByCategory(ctx context.Context, arg ListProductsByCategoryParams) ([]ListProductsByCategoryRow, error) {
	rows, err := q.db.Query(ctx, ListProductsByCategory,
		arg.Column1,
		arg.Status,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListProductsByCategoryRow
	for rows.Next() {
		var i ListProductsByCategoryRow
		if err := rows.Scan(
			&i.ID,
			&i.MerchantID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.Status,
			&i.CategoryID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Images,
			&i.Attributes,
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

const ListRandomProducts = `-- name: ListRandomProducts :many

SELECT p.id,
       p.merchant_id,
       p.name,
       p.description,
       p.price,
       p.status,
       p.category_id,
       p.created_at,
       p.updated_at,
       -- 图片信息
       (SELECT jsonb_agg(jsonb_build_object(
               'url', pi.url,
               'is_primary', pi.is_primary,
               'sort_order', pi.sort_order
                         ))
        FROM products.product_images pi
        WHERE pi.product_id = p.id
          AND pi.merchant_id = p.merchant_id) AS images,
       -- 属性信息
       pa.attributes
FROM products.products p
         LEFT JOIN products.product_attributes pa
                   ON p.id = pa.product_id AND p.merchant_id = pa.merchant_id
WHERE p.status = $1
  AND p.deleted_at IS NULL
ORDER BY random()
LIMIT $2 OFFSET $3
`

type ListRandomProductsParams struct {
	Status int16 `json:"status"`
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

type ListRandomProductsRow struct {
	ID          uuid.UUID      `json:"id"`
	MerchantID  uuid.UUID      `json:"merchantID"`
	Name        string         `json:"name"`
	Description *string        `json:"description"`
	Price       pgtype.Numeric `json:"price"`
	Status      int16          `json:"status"`
	CategoryID  int64          `json:"categoryID"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	Images      []byte         `json:"images"`
	Attributes  []byte         `json:"attributes"`
}

// 实现随机商品列表查询（使用PostgreSQL的TABLESAMPLE优化性能）,数据量太少不显示
// SELECT
//
//	p.id,
//	p.merchant_id,
//	p.name,
//	p.description,
//	p.price,
//	p.status,
//	p.category_id,
//	p.created_at,
//	p.updated_at,
//	-- 图片信息
//	(
//	    SELECT jsonb_agg(jsonb_build_object(
//	            'url', pi.url,
//	            'is_primary', pi.is_primary,
//	            'sort_order', pi.sort_order
//	                     ))
//	    FROM products.product_images pi
//	    WHERE pi.product_id = p.id AND pi.merchant_id = p.merchant_id
//	) AS images,
//	-- 属性信息
//	pa.attributes
//
// FROM products.products p
//
//	TABLESAMPLE BERNOULLI (0.1) REPEATABLE (123)
//	LEFT JOIN products.product_attributes pa
//	          ON p.id = pa.product_id AND p.merchant_id = pa.merchant_id
//
// WHERE p.status = $1 AND p.deleted_at IS NULL
// ORDER BY random()
//
//	   LIMIT $2 OFFSET $3;
//
//
//	SELECT p.id,
//	       p.merchant_id,
//	       p.name,
//	       p.description,
//	       p.price,
//	       p.status,
//	       p.category_id,
//	       p.created_at,
//	       p.updated_at,
//	       -- 图片信息
//	       (SELECT jsonb_agg(jsonb_build_object(
//	               'url', pi.url,
//	               'is_primary', pi.is_primary,
//	               'sort_order', pi.sort_order
//	                         ))
//	        FROM products.product_images pi
//	        WHERE pi.product_id = p.id
//	          AND pi.merchant_id = p.merchant_id) AS images,
//	       -- 属性信息
//	       pa.attributes
//	FROM products.products p
//	         LEFT JOIN products.product_attributes pa
//	                   ON p.id = pa.product_id AND p.merchant_id = pa.merchant_id
//	WHERE p.status = $1
//	  AND p.deleted_at IS NULL
//	ORDER BY random()
//	LIMIT $2 OFFSET $3
func (q *Queries) ListRandomProducts(ctx context.Context, arg ListRandomProductsParams) ([]ListRandomProductsRow, error) {
	rows, err := q.db.Query(ctx, ListRandomProducts, arg.Status, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListRandomProductsRow
	for rows.Next() {
		var i ListRandomProductsRow
		if err := rows.Scan(
			&i.ID,
			&i.MerchantID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.Status,
			&i.CategoryID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Images,
			&i.Attributes,
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

const SearchFullProductsByName = `-- name: SearchFullProductsByName :many
SELECT p.id,
       p.name,
       p.description,
       p.price,
       p.status,
       p.merchant_id,
       p.created_at,
       p.updated_at,
       i.stock,
       (SELECT jsonb_agg(jsonb_build_object(
               'url', pi.url,
               'is_primary', pi.is_primary,
               'sort_order', pi.sort_order
                         ))
        FROM products.product_images pi
        WHERE pi.product_id = p.id
          AND pi.merchant_id = p.merchant_id) AS images,
       pa.attributes
FROM products.products p
         INNER JOIN products.inventory i
                    ON p.id = i.product_id AND p.merchant_id = i.merchant_id
         LEFT JOIN products.product_attributes pa
                   ON p.id = pa.product_id AND p.merchant_id = pa.merchant_id
WHERE p.name ILIKE '%' || $1 || '%'
  AND p.deleted_at IS NULL
ORDER BY ts_rank(to_tsvector('simple', p.name), plainto_tsquery('simple', $2)) DESC,
         p.created_at DESC
LIMIT $4 OFFSET $3
`

type SearchFullProductsByNameParams struct {
	Name     *string `json:"name"`
	Query    string  `json:"query"`
	PageSize int64   `json:"pageSize"`
	Page     int64   `json:"page"`
}

type SearchFullProductsByNameRow struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Description *string        `json:"description"`
	Price       pgtype.Numeric `json:"price"`
	Status      int16          `json:"status"`
	MerchantID  uuid.UUID      `json:"merchantID"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	Stock       int32          `json:"stock"`
	Images      []byte         `json:"images"`
	Attributes  []byte         `json:"attributes"`
}

// 商品搜索查询
//
//	SELECT p.id,
//	       p.name,
//	       p.description,
//	       p.price,
//	       p.status,
//	       p.merchant_id,
//	       p.created_at,
//	       p.updated_at,
//	       i.stock,
//	       (SELECT jsonb_agg(jsonb_build_object(
//	               'url', pi.url,
//	               'is_primary', pi.is_primary,
//	               'sort_order', pi.sort_order
//	                         ))
//	        FROM products.product_images pi
//	        WHERE pi.product_id = p.id
//	          AND pi.merchant_id = p.merchant_id) AS images,
//	       pa.attributes
//	FROM products.products p
//	         INNER JOIN products.inventory i
//	                    ON p.id = i.product_id AND p.merchant_id = i.merchant_id
//	         LEFT JOIN products.product_attributes pa
//	                   ON p.id = pa.product_id AND p.merchant_id = pa.merchant_id
//	WHERE p.name ILIKE '%' || $1 || '%'
//	  AND p.deleted_at IS NULL
//	ORDER BY ts_rank(to_tsvector('simple', p.name), plainto_tsquery('simple', $2)) DESC,
//	         p.created_at DESC
//	LIMIT $4 OFFSET $3
func (q *Queries) SearchFullProductsByName(ctx context.Context, arg SearchFullProductsByNameParams) ([]SearchFullProductsByNameRow, error) {
	rows, err := q.db.Query(ctx, SearchFullProductsByName,
		arg.Name,
		arg.Query,
		arg.PageSize,
		arg.Page,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchFullProductsByNameRow
	for rows.Next() {
		var i SearchFullProductsByNameRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.Status,
			&i.MerchantID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Stock,
			&i.Images,
			&i.Attributes,
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
SET deleted_at = NOW(),
    status     = $3
WHERE merchant_id = $1
  AND id = $2
`

type SoftDeleteProductParams struct {
	MerchantID uuid.UUID `json:"merchantID"`
	ID         uuid.UUID `json:"id"`
	Status     int16     `json:"status"`
}

// 软删除商品，设置删除时间戳
//
//	UPDATE products.products
//	SET deleted_at = NOW(),
//	    status     = $3
//	WHERE merchant_id = $1
//	  AND id = $2
func (q *Queries) SoftDeleteProduct(ctx context.Context, arg SoftDeleteProductParams) error {
	_, err := q.db.Exec(ctx, SoftDeleteProduct, arg.MerchantID, arg.ID, arg.Status)
	return err
}

const UpdateProduct = `-- name: UpdateProduct :exec
UPDATE products.products
SET name        = $2,
    description = $3,
    price       = $4,
    status      = $5,
    updated_at  = NOW()
WHERE id = $1
  AND merchant_id = $6
  AND updated_at = $7
`

type UpdateProductParams struct {
	ID          uuid.UUID          `json:"id"`
	Name        string             `json:"name"`
	Description *string            `json:"description"`
	Price       pgtype.Numeric     `json:"price"`
	Status      int16              `json:"status"`
	MerchantID  uuid.UUID          `json:"merchantID"`
	UpdatedAt   pgtype.Timestamptz `json:"updatedAt"`
}

// 更新商品基础信息，使用乐观锁控制并发
//
//	UPDATE products.products
//	SET name        = $2,
//	    description = $3,
//	    price       = $4,
//	    status      = $5,
//	    updated_at  = NOW()
//	WHERE id = $1
//	  AND merchant_id = $6
//	  AND updated_at = $7
func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) error {
	_, err := q.db.Exec(ctx, UpdateProduct,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.Price,
		arg.Status,
		arg.MerchantID,
		arg.UpdatedAt,
	)
	return err
}

const UpdateProductAttribute = `-- name: UpdateProductAttribute :exec
UPDATE products.product_attributes
SET attributes = $3,
    updated_at = NOW()
WHERE merchant_id = $1
  AND product_id = $2
RETURNING updated_at
`

type UpdateProductAttributeParams struct {
	MerchantID uuid.UUID `json:"merchantID"`
	ProductID  uuid.UUID `json:"productID"`
	Attributes []byte    `json:"attributes"`
}

// UpdateProductAttribute
//
//	UPDATE products.product_attributes
//	SET attributes = $3,
//	    updated_at = NOW()
//	WHERE merchant_id = $1
//	  AND product_id = $2
//	RETURNING updated_at
func (q *Queries) UpdateProductAttribute(ctx context.Context, arg UpdateProductAttributeParams) error {
	_, err := q.db.Exec(ctx, UpdateProductAttribute, arg.MerchantID, arg.ProductID, arg.Attributes)
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
	ID             uuid.UUID   `json:"id"`
	Status         int16       `json:"status"`
	CurrentAuditID pgtype.UUID `json:"currentAuditID"`
	MerchantID     uuid.UUID   `json:"merchantID"`
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
