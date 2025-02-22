// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package models

import (
	"context"
)

const CheckCartItem = `-- name: CheckCartItem :exec

UPDATE cart_schema.cart_items AS ci
SET selected = TRUE
WHERE ci.cart_id = 
    (SELECT c.cart_id
     FROM cart_schema.cart AS c
     WHERE c.user_id = $1 AND c.cart_name = $2 LIMIT 1) 
    AND ci.merchant_id = $3  -- 商家ID
    AND ci.product_id = $4
`

type CheckCartItemParams struct {
	UserID     string `json:"userID"`
	CartName   string `json:"cartName"`
	MerchantID string `json:"merchantID"`
	ProductID  int32  `json:"productID"`
}

// 获取用户的购物车ID
//
//	UPDATE cart_schema.cart_items AS ci
//	SET selected = TRUE
//	WHERE ci.cart_id =
//	    (SELECT c.cart_id
//	     FROM cart_schema.cart AS c
//	     WHERE c.user_id = $1 AND c.cart_name = $2 LIMIT 1)
//	    AND ci.merchant_id = $3  -- 商家ID
//	    AND ci.product_id = $4
func (q *Queries) CheckCartItem(ctx context.Context, arg CheckCartItemParams) error {
	_, err := q.db.Exec(ctx, CheckCartItem,
		arg.UserID,
		arg.CartName,
		arg.MerchantID,
		arg.ProductID,
	)
	return err
}

const CreateCart = `-- name: CreateCart :one
INSERT INTO cart_schema.cart (user_id, cart_name)
VALUES ($1, $2)
RETURNING cart_id, user_id, cart_name, status, created_at, updated_at
`

type CreateCartParams struct {
	UserID   string `json:"userID"`
	CartName string `json:"cartName"`
}

// CreateCart
//
//	INSERT INTO cart_schema.cart (user_id, cart_name)
//	VALUES ($1, $2)
//	RETURNING cart_id, user_id, cart_name, status, created_at, updated_at
func (q *Queries) CreateCart(ctx context.Context, arg CreateCartParams) (CartSchemaCart, error) {
	row := q.db.QueryRow(ctx, CreateCart, arg.UserID, arg.CartName)
	var i CartSchemaCart
	err := row.Scan(
		&i.CartID,
		&i.UserID,
		&i.CartName,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const CreateOrder = `-- name: CreateOrder :many
SELECT ci.merchant_id, ci.product_id, ci.quantity, ci.selected
FROM cart_schema.cart_items AS ci
WHERE ci.cart_id = 
    (SELECT c.cart_id
     FROM cart_schema.cart AS c
     WHERE c.user_id = $1 AND c.cart_name = $2 LIMIT 1) 
    AND ci.selected = TRUE
`

type CreateOrderParams struct {
	UserID   string `json:"userID"`
	CartName string `json:"cartName"`
}

type CreateOrderRow struct {
	MerchantID string `json:"merchantID"`
	ProductID  int32  `json:"productID"`
	Quantity   int32  `json:"quantity"`
	Selected   bool   `json:"selected"`
}

// CreateOrder
//
//	SELECT ci.merchant_id, ci.product_id, ci.quantity, ci.selected
//	FROM cart_schema.cart_items AS ci
//	WHERE ci.cart_id =
//	    (SELECT c.cart_id
//	     FROM cart_schema.cart AS c
//	     WHERE c.user_id = $1 AND c.cart_name = $2 LIMIT 1)
//	    AND ci.selected = TRUE
func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) ([]CreateOrderRow, error) {
	rows, err := q.db.Query(ctx, CreateOrder, arg.UserID, arg.CartName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CreateOrderRow
	for rows.Next() {
		var i CreateOrderRow
		if err := rows.Scan(
			&i.MerchantID,
			&i.ProductID,
			&i.Quantity,
			&i.Selected,
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

const EmptyCart = `-- name: EmptyCart :exec
DELETE FROM cart_schema.cart_items AS ci
WHERE ci.cart_id = 
    (SELECT c.cart_id
     FROM cart_schema.cart AS c
     WHERE c.user_id = $1 AND c.cart_name = $2)
`

type EmptyCartParams struct {
	UserID   string `json:"userID"`
	CartName string `json:"cartName"`
}

// EmptyCart
//
//	DELETE FROM cart_schema.cart_items AS ci
//	WHERE ci.cart_id =
//	    (SELECT c.cart_id
//	     FROM cart_schema.cart AS c
//	     WHERE c.user_id = $1 AND c.cart_name = $2)
func (q *Queries) EmptyCart(ctx context.Context, arg EmptyCartParams) error {
	_, err := q.db.Exec(ctx, EmptyCart, arg.UserID, arg.CartName)
	return err
}

const GetCart = `-- name: GetCart :many
SELECT ci.merchant_id, ci.product_id, ci.quantity, ci.selected 
FROM cart_schema.cart_items AS ci
WHERE ci.cart_id = 
    (SELECT c.cart_id
     FROM cart_schema.cart AS c
     WHERE c.user_id = $1 AND c.cart_name = $2 LIMIT 1)
`

type GetCartParams struct {
	UserID   string `json:"userID"`
	CartName string `json:"cartName"`
}

type GetCartRow struct {
	MerchantID string `json:"merchantID"`
	ProductID  int32  `json:"productID"`
	Quantity   int32  `json:"quantity"`
	Selected   bool   `json:"selected"`
}

// GetCart
//
//	SELECT ci.merchant_id, ci.product_id, ci.quantity, ci.selected
//	FROM cart_schema.cart_items AS ci
//	WHERE ci.cart_id =
//	    (SELECT c.cart_id
//	     FROM cart_schema.cart AS c
//	     WHERE c.user_id = $1 AND c.cart_name = $2 LIMIT 1)
func (q *Queries) GetCart(ctx context.Context, arg GetCartParams) ([]GetCartRow, error) {
	rows, err := q.db.Query(ctx, GetCart, arg.UserID, arg.CartName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCartRow
	for rows.Next() {
		var i GetCartRow
		if err := rows.Scan(
			&i.MerchantID,
			&i.ProductID,
			&i.Quantity,
			&i.Selected,
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

const ListCarts = `-- name: ListCarts :many
SELECT c.cart_id, c.cart_name
FROM cart_schema.cart AS c
WHERE c.user_id = $1
`

type ListCartsRow struct {
	CartID   int32  `json:"cartID"`
	CartName string `json:"cartName"`
}

// ListCarts
//
//	SELECT c.cart_id, c.cart_name
//	FROM cart_schema.cart AS c
//	WHERE c.user_id = $1
func (q *Queries) ListCarts(ctx context.Context, userID string) ([]ListCartsRow, error) {
	rows, err := q.db.Query(ctx, ListCarts, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListCartsRow
	for rows.Next() {
		var i ListCartsRow
		if err := rows.Scan(&i.CartID, &i.CartName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const RemoveCartItem = `-- name: RemoveCartItem :one

DELETE FROM cart_schema.cart_items AS ci
WHERE ci.cart_id = 
    (SELECT c.cart_id
     FROM cart_schema.cart AS c
     WHERE c.user_id = $1 AND c.cart_name = $2 LIMIT 1)  -- 获取用户的购物车ID
    AND ci.merchant_id = $3  -- 商家ID
    AND ci.product_id = $4  -- 删除指定商品ID
RETURNING cart_item_id, cart_id, merchant_id, product_id, quantity, selected, created_at, updated_at
`

type RemoveCartItemParams struct {
	UserID     string `json:"userID"`
	CartName   string `json:"cartName"`
	MerchantID string `json:"merchantID"`
	ProductID  int32  `json:"productID"`
}

// 获取用户的购物车ID
//
//	DELETE FROM cart_schema.cart_items AS ci
//	WHERE ci.cart_id =
//	    (SELECT c.cart_id
//	     FROM cart_schema.cart AS c
//	     WHERE c.user_id = $1 AND c.cart_name = $2 LIMIT 1)  -- 获取用户的购物车ID
//	    AND ci.merchant_id = $3  -- 商家ID
//	    AND ci.product_id = $4  -- 删除指定商品ID
//	RETURNING cart_item_id, cart_id, merchant_id, product_id, quantity, selected, created_at, updated_at
func (q *Queries) RemoveCartItem(ctx context.Context, arg RemoveCartItemParams) (CartSchemaCartItems, error) {
	row := q.db.QueryRow(ctx, RemoveCartItem,
		arg.UserID,
		arg.CartName,
		arg.MerchantID,
		arg.ProductID,
	)
	var i CartSchemaCartItems
	err := row.Scan(
		&i.CartItemID,
		&i.CartID,
		&i.MerchantID,
		&i.ProductID,
		&i.Quantity,
		&i.Selected,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const UncheckCartItem = `-- name: UncheckCartItem :exec
UPDATE cart_schema.cart_items AS ci
SET selected = FALSE
WHERE ci.cart_id = 
    (SELECT c.cart_id
     FROM cart_schema.cart AS c
     WHERE c.user_id = $1 AND c.cart_name = $2 LIMIT 1) 
    AND ci.merchant_id = $3  -- 商家ID
    AND ci.product_id = $4
`

type UncheckCartItemParams struct {
	UserID     string `json:"userID"`
	CartName   string `json:"cartName"`
	MerchantID string `json:"merchantID"`
	ProductID  int32  `json:"productID"`
}

// UncheckCartItem
//
//	UPDATE cart_schema.cart_items AS ci
//	SET selected = FALSE
//	WHERE ci.cart_id =
//	    (SELECT c.cart_id
//	     FROM cart_schema.cart AS c
//	     WHERE c.user_id = $1 AND c.cart_name = $2 LIMIT 1)
//	    AND ci.merchant_id = $3  -- 商家ID
//	    AND ci.product_id = $4
func (q *Queries) UncheckCartItem(ctx context.Context, arg UncheckCartItemParams) error {
	_, err := q.db.Exec(ctx, UncheckCartItem,
		arg.UserID,
		arg.CartName,
		arg.MerchantID,
		arg.ProductID,
	)
	return err
}

const UpsertItem = `-- name: UpsertItem :one
WITH cart_id_cte AS (
    SELECT c.cart_id
    FROM cart_schema.cart AS c
    WHERE c.user_id = $1 AND c.cart_name = $2
    LIMIT 1
),
insert_cart AS (
    INSERT INTO cart_schema.cart (user_id, cart_name)
    SELECT $1, $2
    WHERE NOT EXISTS (SELECT 1 FROM cart_id_cte)
    RETURNING cart_id
)
INSERT INTO cart_schema.cart_items (cart_id, merchant_id, product_id, quantity, created_at, updated_at)
VALUES (
    COALESCE((SELECT cart_id FROM cart_id_cte), (SELECT cart_id FROM insert_cart)),  -- 获取或创建购物车ID
    $3,   -- 商家ID
    $4,   -- 商品ID
    $5,   -- 商品数量
    CURRENT_TIMESTAMP,  -- 创建时间
    CURRENT_TIMESTAMP   -- 更新时间
)
ON CONFLICT (cart_id, merchant_id, product_id)  -- 如果购物车ID、商家ID和商品ID组合重复
DO UPDATE SET 
    quantity = EXCLUDED.quantity,  -- 更新商品数量
    updated_at = CURRENT_TIMESTAMP  -- 更新时间
RETURNING cart_item_id, cart_id, merchant_id, product_id, quantity, selected, created_at, updated_at
`

type UpsertItemParams struct {
	UserID     string `json:"userID"`
	CartName   string `json:"cartName"`
	MerchantID string `json:"merchantID"`
	ProductID  int32  `json:"productID"`
	Quantity   int32  `json:"quantity"`
}

// UpsertItem
//
//	WITH cart_id_cte AS (
//	    SELECT c.cart_id
//	    FROM cart_schema.cart AS c
//	    WHERE c.user_id = $1 AND c.cart_name = $2
//	    LIMIT 1
//	),
//	insert_cart AS (
//	    INSERT INTO cart_schema.cart (user_id, cart_name)
//	    SELECT $1, $2
//	    WHERE NOT EXISTS (SELECT 1 FROM cart_id_cte)
//	    RETURNING cart_id
//	)
//	INSERT INTO cart_schema.cart_items (cart_id, merchant_id, product_id, quantity, created_at, updated_at)
//	VALUES (
//	    COALESCE((SELECT cart_id FROM cart_id_cte), (SELECT cart_id FROM insert_cart)),  -- 获取或创建购物车ID
//	    $3,   -- 商家ID
//	    $4,   -- 商品ID
//	    $5,   -- 商品数量
//	    CURRENT_TIMESTAMP,  -- 创建时间
//	    CURRENT_TIMESTAMP   -- 更新时间
//	)
//	ON CONFLICT (cart_id, merchant_id, product_id)  -- 如果购物车ID、商家ID和商品ID组合重复
//	DO UPDATE SET
//	    quantity = EXCLUDED.quantity,  -- 更新商品数量
//	    updated_at = CURRENT_TIMESTAMP  -- 更新时间
//	RETURNING cart_item_id, cart_id, merchant_id, product_id, quantity, selected, created_at, updated_at
func (q *Queries) UpsertItem(ctx context.Context, arg UpsertItemParams) (CartSchemaCartItems, error) {
	row := q.db.QueryRow(ctx, UpsertItem,
		arg.UserID,
		arg.CartName,
		arg.MerchantID,
		arg.ProductID,
		arg.Quantity,
	)
	var i CartSchemaCartItems
	err := row.Scan(
		&i.CartItemID,
		&i.CartID,
		&i.MerchantID,
		&i.ProductID,
		&i.Quantity,
		&i.Selected,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
