// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: account.sql

package models

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const DeleteFavorites = `-- name: DeleteFavorites :exec
DELETE
FROM users.favorites
WHERE user_id = $1
  AND product_id = $2
  AND merchant_id = $3
`

type DeleteFavoritesParams struct {
	UserID     uuid.UUID `json:"userID"`
	ProductID  uuid.UUID `json:"productID"`
	MerchantID uuid.UUID `json:"merchantID"`
}

// DeleteFavorites
//
//	DELETE
//	FROM users.favorites
//	WHERE user_id = $1
//	  AND product_id = $2
//	  AND merchant_id = $3
func (q *Queries) DeleteFavorites(ctx context.Context, arg DeleteFavoritesParams) error {
	_, err := q.db.Exec(ctx, DeleteFavorites, arg.UserID, arg.ProductID, arg.MerchantID)
	return err
}

const GetFavorites = `-- name: GetFavorites :many
SELECT p.id,
       p.merchant_id,
       p.name,
       p.description,
       p.price,
       p.status,
       p.category_id,
       p.created_at,
       p.updated_at,
       i.stock,
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
         INNER JOIN products.inventory i
                    ON p.id = i.product_id AND p.merchant_id = i.merchant_id
         LEFT JOIN products.product_attributes pa
                   ON p.id = pa.product_id AND p.merchant_id = pa.merchant_id
         JOIN users.favorites uf ON p.id = uf.product_id
WHERE uf.user_id = $1::UUID
ORDER BY uf.created_at DESC
LIMIT $3::int OFFSET $2::int
`

type GetFavoritesParams struct {
	UserID   pgtype.UUID `json:"userID"`
	Page     *int32      `json:"page"`
	PageSize *int32      `json:"pageSize"`
}

type GetFavoritesRow struct {
	ID          uuid.UUID          `json:"id"`
	MerchantID  uuid.UUID          `json:"merchantID"`
	Name        string             `json:"name"`
	Description *string            `json:"description"`
	Price       interface{}        `json:"price"`
	Status      int16              `json:"status"`
	CategoryID  int64              `json:"categoryID"`
	CreatedAt   pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt   pgtype.Timestamptz `json:"updatedAt"`
	Stock       int32              `json:"stock"`
	Images      []byte             `json:"images"`
	Attributes  []byte             `json:"attributes"`
}

// GetFavorites
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
//	       i.stock,
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
//	         INNER JOIN products.inventory i
//	                    ON p.id = i.product_id AND p.merchant_id = i.merchant_id
//	         LEFT JOIN products.product_attributes pa
//	                   ON p.id = pa.product_id AND p.merchant_id = pa.merchant_id
//	         JOIN users.favorites uf ON p.id = uf.product_id
//	WHERE uf.user_id = $1::UUID
//	ORDER BY uf.created_at DESC
//	LIMIT $3::int OFFSET $2::int
func (q *Queries) GetFavorites(ctx context.Context, arg GetFavoritesParams) ([]GetFavoritesRow, error) {
	rows, err := q.db.Query(ctx, GetFavorites, arg.UserID, arg.Page, arg.PageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFavoritesRow
	for rows.Next() {
		var i GetFavoritesRow
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

const SetFavorites = `-- name: SetFavorites :exec
INSERT INTO users.favorites(user_id, product_id, merchant_id)
VALUES ($1, $2, $3)
`

type SetFavoritesParams struct {
	UserID     uuid.UUID `json:"userID"`
	ProductID  uuid.UUID `json:"productID"`
	MerchantID uuid.UUID `json:"merchantID"`
}

// SetFavorites
//
//	INSERT INTO users.favorites(user_id, product_id, merchant_id)
//	VALUES ($1, $2, $3)
func (q *Queries) SetFavorites(ctx context.Context, arg SetFavoritesParams) error {
	_, err := q.db.Exec(ctx, SetFavorites, arg.UserID, arg.ProductID, arg.MerchantID)
	return err
}
