// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package models

import (
	"context"
)

const DeleteProduct = `-- name: DeleteProduct :exec
DELETE FROM products.products
WHERE id = $1
`

type DeleteProductParams struct {
	ID int32 `json:"id"`
}

//DeleteProduct
//
//  DELETE FROM products.products
//  WHERE id = $1
func (q *Queries) DeleteProduct(ctx context.Context, arg DeleteProductParams)(ProductsProducts, error) {
	row := q.db.QueryRow(ctx, DeleteProduct, arg.ID)
	var i ProductsProducts
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Picture,
		&i.Price,
		&i.Categories,
	)
	return i, err
}

const UpdateProduct = `-- name: UpdateProduct :exec
UPDATE products.products
SET name = $2, description = $3, picture = $4, price = $5, categories = $6
WHERE id = $1
`

type UpdateProductParams struct {
	ID          int32    `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Picture     string   `json:"picture"`
	Price       float32  `json:"price"`
	Categories  []string `json:"categories"`
}

//UpdateProduct

//  UPDATE products.products
//  SET name = $2, description = $3, picture = $4, price = $5, categories = $6
//  WHERE id = $1
func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams)(ProductsProducts, error) {
	row:= q.db.QueryRow(ctx, UpdateProduct, 
		arg.ID, 
		arg.Name,
		arg.Description, 
		arg.Picture, 
		arg.Price, 
		arg.Categories,
	)
	var i ProductsProducts
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Picture,
		&i.Price,
		&i.Categories,
	)
	return i, err
}

const CreateProduct = `-- name: CreateProduct :one
INSERT INTO products.products(name, description, picture, price, categories)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, description, picture, price, categories
`

type CreateProductParams struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Picture     string   `json:"picture"`
	Price       float32  `json:"price"`
	Categories  []string `json:"categories"`
}

//CreateProduct
//
//  INSERT INTO products.products(name, description, picture, price, categories)
//  VALUES ($1, $2, $3, $4, $5)
//  RETURNING id, name, description, picture, price, categories
func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (ProductsProducts, error) {
	row := q.db.QueryRow(ctx, CreateProduct,
		arg.Name,
		arg.Description,
		arg.Picture,
		arg.Price,
		arg.Categories,
	)
	var i ProductsProducts
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Picture,
		&i.Price,
		&i.Categories,
	)
	return i, err
}

const GetProduct = `-- name: GetProduct :one
SELECT id, name, description, picture, price, categories
FROM products.products
WHERE id = $1
LIMIT 1
`

//GetProduct
//
//  SELECT id, name, description, picture, price, categories
//  FROM products.products
//  WHERE id = $1
//  LIMIT 1
func (q *Queries) GetProduct(ctx context.Context, id int32) (ProductsProducts, error) {
	row := q.db.QueryRow(ctx, GetProduct, id)
	var i ProductsProducts
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Picture,
		&i.Price,
		&i.Categories,
	)
	return i, err
}

const ListProducts = `-- name: ListProducts :many
SELECT id, name, description, picture, price, categories
FROM products.products
ORDER BY id
OFFSET $1 LIMIT $2
`

type ListProductsParams struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"pageSize"`
}

//ListProducts
//
//  SELECT id, name, description, picture, price, categories
//  FROM products.products
//  ORDER BY id
//  OFFSET $1 LIMIT $2
func (q *Queries) ListProducts(ctx context.Context, arg ListProductsParams) ([]ProductsProducts, error) {
	rows, err := q.db.Query(ctx, ListProducts, arg.Page, arg.PageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ProductsProducts
	for rows.Next() {
		var i ProductsProducts
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Picture,
			&i.Price,
			&i.Categories,
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

const SearchProducts = `-- name: SearchProducts :many
SELECT id, name, description, picture, price, categories
FROM products.products
WHERE name ILIKE '%' || $1 || '%'
`

//SearchProducts
//
//  SELECT id, name, description, picture, price, categories
//  FROM products.products
//  WHERE name ILIKE '%' || $1 || '%'
func (q *Queries) SearchProducts(ctx context.Context, name *string) ([]ProductsProducts, error) {
	rows, err := q.db.Query(ctx, SearchProducts, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ProductsProducts
	for rows.Next() {
		var i ProductsProducts
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Picture,
			&i.Price,
			&i.Categories,
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
