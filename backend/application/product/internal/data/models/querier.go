// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package models

import (
	"context"
)

type Querier interface {
	//CreateProduct
	//
	//  INSERT INTO products.products(name, description, picture, price, categories)
	//  VALUES ($1, $2, $3, $4, $5)
	//  RETURNING id, name, description, picture, price, categories
	CreateProduct(ctx context.Context, arg CreateProductParams) (ProductsProducts, error)
	//GetProduct
	//
	//  SELECT id, name, description, picture, price, categories
	//  FROM products.products
	//  WHERE id = $1
	//  LIMIT 1
	GetProduct(ctx context.Context, id int32) (ProductsProducts, error)
	//ListProducts
	//
	//  SELECT id, name, description, picture, price, categories
	//  FROM products.products
	//  ORDER BY id
	//  OFFSET $1 LIMIT $2
	ListProducts(ctx context.Context, arg ListProductsParams) ([]ProductsProducts, error)
	//SearchProducts
	//
	//  SELECT id, name, description, picture, price, categories
	//  FROM products.products
	//  WHERE name ILIKE '%' || $1 || '%'
	SearchProducts(ctx context.Context, name *string) ([]ProductsProducts, error)
}

var _ Querier = (*Queries)(nil)
