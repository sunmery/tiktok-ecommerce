// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package models

import (
	"context"
)

const EmptyCart = `-- name: EmptyCart :many
DELETE FROM cart_schema.cart_items AS ci
WHERE ci.cart_id = 
    (SELECT c.cart_id
     FROM cart_schema.cart AS c
     WHERE c.user_id = $1)  -- 获取用户的购物车ID
RETURNING cart_item_id, cart_id, product_id, quantity, created_at, updated_at
`

// EmptyCart
//
//	DELETE FROM cart_schema.cart_items AS ci
//	WHERE ci.cart_id =
//	    (SELECT c.cart_id
//	     FROM cart_schema.cart AS c
//	     WHERE c.user_id = $1)  -- 获取用户的购物车ID
//	RETURNING cart_item_id, cart_id, product_id, quantity, created_at, updated_at
func (q *Queries) EmptyCart(ctx context.Context, userID int32) ([]CartSchemaCartItems, error) {
	rows, err := q.db.Query(ctx, EmptyCart, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CartSchemaCartItems
	for rows.Next() {
		var i CartSchemaCartItems
		if err := rows.Scan(
			&i.CartItemID,
			&i.CartID,
			&i.ProductID,
			&i.Quantity,
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

const GetCart = `-- name: GetCart :many
SELECT ci.cart_item_id, ci.quantity 
FROM cart_schema.cart_items AS ci
WHERE ci.cart_id = 
    (SELECT c.cart_id
     FROM cart_schema.cart AS c
     WHERE c.user_id = $1)
`

type GetCartRow struct {
	CartItemID int32 `json:"cartItemID"`
	Quantity   int32 `json:"quantity"`
}

// GetCart
//
//	SELECT ci.cart_item_id, ci.quantity
//	FROM cart_schema.cart_items AS ci
//	WHERE ci.cart_id =
//	    (SELECT c.cart_id
//	     FROM cart_schema.cart AS c
//	     WHERE c.user_id = $1)
func (q *Queries) GetCart(ctx context.Context, userID int32) ([]GetCartRow, error) {
	rows, err := q.db.Query(ctx, GetCart, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCartRow
	for rows.Next() {
		var i GetCartRow
		if err := rows.Scan(&i.CartItemID, &i.Quantity); err != nil {
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
     WHERE c.user_id = $1)  -- 获取用户的购物车ID
    AND ci.product_id = $2  -- 删除指定商品ID
RETURNING cart_item_id, cart_id, product_id, quantity, created_at, updated_at
`

type RemoveCartItemParams struct {
	UserID    int32 `json:"userID"`
	ProductID int32 `json:"productID"`
}

// 获取用户的购物车ID
//
//	DELETE FROM cart_schema.cart_items AS ci
//	WHERE ci.cart_id =
//	    (SELECT c.cart_id
//	     FROM cart_schema.cart AS c
//	     WHERE c.user_id = $1)  -- 获取用户的购物车ID
//	    AND ci.product_id = $2  -- 删除指定商品ID
//	RETURNING cart_item_id, cart_id, product_id, quantity, created_at, updated_at
func (q *Queries) RemoveCartItem(ctx context.Context, arg RemoveCartItemParams) (CartSchemaCartItems, error) {
	row := q.db.QueryRow(ctx, RemoveCartItem, arg.UserID, arg.ProductID)
	var i CartSchemaCartItems
	err := row.Scan(
		&i.CartItemID,
		&i.CartID,
		&i.ProductID,
		&i.Quantity,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const UpsertItem = `-- name: UpsertItem :one
INSERT INTO cart_schema.cart_items (cart_id, product_id, quantity, created_at, updated_at)
VALUES (
    (SELECT c.cart_id
     FROM cart_schema.cart AS c
     WHERE c.user_id = $1 LIMIT 1),  -- 获取用户的购物车ID
    $2,   -- 商品ID
    $3,   -- 商品数量
    CURRENT_TIMESTAMP,  -- 创建时间
    CURRENT_TIMESTAMP   -- 更新时间
)
ON CONFLICT (cart_id, product_id)  -- 如果购物车ID和商品ID组合重复
DO UPDATE SET 
    quantity = cart_schema.cart_items.quantity + EXCLUDED.quantity,  -- 更新商品数量
    updated_at = CURRENT_TIMESTAMP  -- 更新时间
RETURNING cart_item_id, cart_id, product_id, quantity, created_at, updated_at
`

type UpsertItemParams struct {
	UserID    int32 `json:"userID"`
	ProductID int32 `json:"productID"`
	Quantity  int32 `json:"quantity"`
}

// UpsertItem
//
//	INSERT INTO cart_schema.cart_items (cart_id, product_id, quantity, created_at, updated_at)
//	VALUES (
//	    (SELECT c.cart_id
//	     FROM cart_schema.cart AS c
//	     WHERE c.user_id = $1 LIMIT 1),  -- 获取用户的购物车ID
//	    $2,   -- 商品ID
//	    $3,   -- 商品数量
//	    CURRENT_TIMESTAMP,  -- 创建时间
//	    CURRENT_TIMESTAMP   -- 更新时间
//	)
//	ON CONFLICT (cart_id, product_id)  -- 如果购物车ID和商品ID组合重复
//	DO UPDATE SET
//	    quantity = cart_schema.cart_items.quantity + EXCLUDED.quantity,  -- 更新商品数量
//	    updated_at = CURRENT_TIMESTAMP  -- 更新时间
//	RETURNING cart_item_id, cart_id, product_id, quantity, created_at, updated_at
func (q *Queries) UpsertItem(ctx context.Context, arg UpsertItemParams) (CartSchemaCartItems, error) {
	row := q.db.QueryRow(ctx, UpsertItem, arg.UserID, arg.ProductID, arg.Quantity)
	var i CartSchemaCartItems
	err := row.Scan(
		&i.CartItemID,
		&i.CartID,
		&i.ProductID,
		&i.Quantity,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
