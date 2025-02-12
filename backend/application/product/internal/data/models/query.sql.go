// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package models

import (
	"context"
)

const CreateAuditLog = `-- name: CreateAuditLog :one
INSERT INTO products.inventory_history(change_reason, product_id, new_stock, owner, username)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, product_id, old_stock, new_stock, change_reason, owner, username, created_at
`

type CreateAuditLogParams struct {
	ChangeReason string `json:"changeReason"`
	ProductID    int32  `json:"productID"`
	NewStock     int32  `json:"newStock"`
	Owner        string `json:"owner"`
	Username     string `json:"username"`
}

// 创建审计日志
//
//  INSERT INTO products.inventory_history(change_reason, product_id, new_stock, owner, username)
//  VALUES ($1, $2, $3, $4, $5)
//  RETURNING id, product_id, old_stock, new_stock, change_reason, owner, username, created_at
func (q *Queries) CreateAuditLog(ctx context.Context, arg CreateAuditLogParams) (ProductsInventoryHistory, error) {
	row := q.db.QueryRow(ctx, CreateAuditLog,
		arg.ChangeReason,
		arg.ProductID,
		arg.NewStock,
		arg.Owner,
		arg.Username,
	)
	var i ProductsInventoryHistory
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.OldStock,
		&i.NewStock,
		&i.ChangeReason,
		&i.Owner,
		&i.Username,
		&i.CreatedAt,
	)
	return i, err
}

const CreateCategories = `-- name: CreateCategories :one
INSERT INTO products.categories (name, parent_id)
VALUES ($1, $2)
RETURNING id, name, parent_id, is_active, created_at
`

type CreateCategoriesParams struct {
	Name     string `json:"name"`
	ParentID *int32 `json:"parentID"`
}

// 创建分类数据
//
//  INSERT INTO products.categories (name, parent_id)
//  VALUES ($1, $2)
//  RETURNING id, name, parent_id, is_active, created_at
func (q *Queries) CreateCategories(ctx context.Context, arg CreateCategoriesParams) (ProductsCategories, error) {
	row := q.db.QueryRow(ctx, CreateCategories, arg.Name, arg.ParentID)
	var i ProductsCategories
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ParentID,
		&i.IsActive,
		&i.CreatedAt,
	)
	return i, err
}

const CreateProduct = `-- name: CreateProduct :one
INSERT INTO products.products(name,
                              description,
                              picture,
                              price,
                              category_id,
                              total_stock)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, name, description, picture, price, category_id, total_stock, available_stock, reserved_stock, low_stock_threshold, allow_negative, created_at, updated_at, version
`

type CreateProductParams struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Picture     string  `json:"picture"`
	Price       float32 `json:"price"`
	CategoryID  []int32 `json:"categoryID"`
	TotalStock  int32   `json:"totalStock"`
}

// 创建商品
//
//  INSERT INTO products.products(name,
//                                description,
//                                picture,
//                                price,
//                                category_id,
//                                total_stock)
//  VALUES ($1, $2, $3, $4, $5, $6)
//  RETURNING id, name, description, picture, price, category_id, total_stock, available_stock, reserved_stock, low_stock_threshold, allow_negative, created_at, updated_at, version
func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (ProductsProducts, error) {
	row := q.db.QueryRow(ctx, CreateProduct,
		arg.Name,
		arg.Description,
		arg.Picture,
		arg.Price,
		arg.CategoryID,
		arg.TotalStock,
	)
	var i ProductsProducts
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Picture,
		&i.Price,
		&i.CategoryID,
		&i.TotalStock,
		&i.AvailableStock,
		&i.ReservedStock,
		&i.LowStockThreshold,
		&i.AllowNegative,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Version,
	)
	return i, err
}

const CreateProductCategories = `-- name: CreateProductCategories :one
INSERT INTO products.product_categories (product_id, category_id)
VALUES ($1, $2)
RETURNING product_id, category_id
`

type CreateProductCategoriesParams struct {
	ProductID  int32 `json:"productID"`
	CategoryID int32 `json:"categoryID"`
}

// 关联商品与分类
// 将商品1关联到分类2（Smartphones）
//
//  INSERT INTO products.product_categories (product_id, category_id)
//  VALUES ($1, $2)
//  RETURNING product_id, category_id
func (q *Queries) CreateProductCategories(ctx context.Context, arg CreateProductCategoriesParams) (ProductsProductCategories, error) {
	row := q.db.QueryRow(ctx, CreateProductCategories, arg.ProductID, arg.CategoryID)
	var i ProductsProductCategories
	err := row.Scan(&i.ProductID, &i.CategoryID)
	return i, err
}

const CreateProductInventoryHistory = `-- name: CreateProductInventoryHistory :one
INSERT INTO products.inventory_history (product_id,
                                        old_stock,
                                        new_stock,
                                        change_reason)
VALUES ($1,
        (SELECT total_stock FROM products.products WHERE id = $1),
        (SELECT total_stock FROM products.products WHERE id = $1) - $2,
        'ORDER_RESERVED')
RETURNING id, product_id, old_stock, new_stock, change_reason, owner, username, created_at
`

type CreateProductInventoryHistoryParams struct {
	Column1 *int32 `json:"column1"`
	Column2 *int32 `json:"column2"`
}

// 记录库存变更
//
//  INSERT INTO products.inventory_history (product_id,
//                                          old_stock,
//                                          new_stock,
//                                          change_reason)
//  VALUES ($1,
//          (SELECT total_stock FROM products.products WHERE id = $1),
//          (SELECT total_stock FROM products.products WHERE id = $1) - $2,
//          'ORDER_RESERVED')
//  RETURNING id, product_id, old_stock, new_stock, change_reason, owner, username, created_at
func (q *Queries) CreateProductInventoryHistory(ctx context.Context, arg CreateProductInventoryHistoryParams) (ProductsInventoryHistory, error) {
	row := q.db.QueryRow(ctx, CreateProductInventoryHistory, arg.Column1, arg.Column2)
	var i ProductsInventoryHistory
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.OldStock,
		&i.NewStock,
		&i.ChangeReason,
		&i.Owner,
		&i.Username,
		&i.CreatedAt,
	)
	return i, err
}

const DeleteProduct = `-- name: DeleteProduct :one
DELETE FROM products.products
WHERE id = $1
RETURNING id, name, description, picture, price, category_id, total_stock, available_stock, reserved_stock, low_stock_threshold, allow_negative, created_at, updated_at, version
`

//DeleteProduct
//
//  DELETE FROM products.products
//  WHERE id = $1
//  RETURNING id, name, description, picture, price, category_id, total_stock, available_stock, reserved_stock, low_stock_threshold, allow_negative, created_at, updated_at, version
func (q *Queries) DeleteProduct(ctx context.Context, id int32) (ProductsProducts, error) {
	row := q.db.QueryRow(ctx, DeleteProduct, id)
	var i ProductsProducts
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Picture,
		&i.Price,
		&i.CategoryID,
		&i.TotalStock,
		&i.AvailableStock,
		&i.ReservedStock,
		&i.LowStockThreshold,
		&i.AllowNegative,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Version,
	)
	return i, err
}

const GetProduct = `-- name: GetProduct :one
SELECT id, name, description, picture, price, category_id, total_stock, available_stock, reserved_stock, low_stock_threshold, allow_negative, created_at, updated_at, version
FROM products.products
WHERE id = $1
LIMIT 1
`

//GetProduct
//
//  SELECT id, name, description, picture, price, category_id, total_stock, available_stock, reserved_stock, low_stock_threshold, allow_negative, created_at, updated_at, version
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
		&i.CategoryID,
		&i.TotalStock,
		&i.AvailableStock,
		&i.ReservedStock,
		&i.LowStockThreshold,
		&i.AllowNegative,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Version,
	)
	return i, err
}

const GetProductCategories = `-- name: GetProductCategories :many
SELECT p.id, p.name, p.description, p.picture, p.price, p.category_id, p.total_stock, p.available_stock, p.reserved_stock, p.low_stock_threshold, p.allow_negative, p.created_at, p.updated_at, p.version
FROM products.products p
         JOIN products.product_categories pc ON p.id = pc.product_id
WHERE pc.category_id = $1
`

// 查询某分类下的所有商品
//
//  SELECT p.id, p.name, p.description, p.picture, p.price, p.category_id, p.total_stock, p.available_stock, p.reserved_stock, p.low_stock_threshold, p.allow_negative, p.created_at, p.updated_at, p.version
//  FROM products.products p
//           JOIN products.product_categories pc ON p.id = pc.product_id
//  WHERE pc.category_id = $1
func (q *Queries) GetProductCategories(ctx context.Context, categoryID int32) ([]ProductsProducts, error) {
	rows, err := q.db.Query(ctx, GetProductCategories, categoryID)
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
			&i.CategoryID,
			&i.TotalStock,
			&i.AvailableStock,
			&i.ReservedStock,
			&i.LowStockThreshold,
			&i.AllowNegative,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Version,
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

const ListProducts = `-- name: ListProducts :many

SELECT id, name, description, picture, price, category_id, total_stock, available_stock, reserved_stock, low_stock_threshold, allow_negative, created_at, updated_at, version
FROM products.products
WHERE ($1 = ANY(category_id))
ORDER BY id
OFFSET $2 LIMIT $3
`

type ListProductsParams struct {
	CategoryID int32 `json:"categoryID"`
	Offset     int64 `json:"offset"`
	Limit      int64 `json:"limit"`
}

//ListProducts
//
//
//  SELECT id, name, description, picture, price, category_id, total_stock, available_stock, reserved_stock, low_stock_threshold, allow_negative, created_at, updated_at, version
//  FROM products.products
//  WHERE ($1 = ANY(category_id))
//  ORDER BY id
//  OFFSET $2 LIMIT $3
func (q *Queries) ListProducts(ctx context.Context, arg ListProductsParams) ([]ProductsProducts, error) {
	rows, err := q.db.Query(ctx, ListProducts, arg.CategoryID, arg.Offset, arg.Limit)
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
			&i.CategoryID,
			&i.TotalStock,
			&i.AvailableStock,
			&i.ReservedStock,
			&i.LowStockThreshold,
			&i.AllowNegative,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Version,
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
SELECT id, name, description, picture, price, category_id, total_stock, available_stock, reserved_stock, low_stock_threshold, allow_negative, created_at, updated_at, version
FROM products.products
WHERE name ILIKE '%' || $1 || '%'
`

//SearchProducts
//
//  SELECT id, name, description, picture, price, category_id, total_stock, available_stock, reserved_stock, low_stock_threshold, allow_negative, created_at, updated_at, version
//  FROM products.products
//  WHERE name ILIKE '%' || $1 || '%'
func (q *Queries) SearchProducts(ctx context.Context, dollar_1 *string) ([]ProductsProducts, error) {
	rows, err := q.db.Query(ctx, SearchProducts, dollar_1)
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
			&i.CategoryID,
			&i.TotalStock,
			&i.AvailableStock,
			&i.ReservedStock,
			&i.LowStockThreshold,
			&i.AllowNegative,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Version,
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

const UpdateAuditLog = `-- name: UpdateAuditLog :one
UPDATE products.inventory_history
SET change_reason = $1, new_stock = $2, owner = $3, username = $4
WHERE product_id = $5
RETURNING id, product_id, old_stock, new_stock, change_reason, owner, username, created_at
`

type UpdateAuditLogParams struct {
	ChangeReason string `json:"changeReason"`
	NewStock     int32  `json:"newStock"`
	Owner        string `json:"owner"`
	Username     string `json:"username"`
	ProductID    int32  `json:"productID"`
}

// 更新审计日志
//
//  UPDATE products.inventory_history
//  SET change_reason = $1, new_stock = $2, owner = $3, username = $4
//  WHERE product_id = $5
//  RETURNING id, product_id, old_stock, new_stock, change_reason, owner, username, created_at
func (q *Queries) UpdateAuditLog(ctx context.Context, arg UpdateAuditLogParams) (ProductsInventoryHistory, error) {
	row := q.db.QueryRow(ctx, UpdateAuditLog,
		arg.ChangeReason,
		arg.NewStock,
		arg.Owner,
		arg.Username,
		arg.ProductID,
	)
	var i ProductsInventoryHistory
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.OldStock,
		&i.NewStock,
		&i.ChangeReason,
		&i.Owner,
		&i.Username,
		&i.CreatedAt,
	)
	return i, err
}

const UpdateProduct = `-- name: UpdateProduct :one
UPDATE products.products
SET name = $1, description = $2, picture = $3, price = $4, category_Id = $5, total_stock = $6
WHERE id = $7
RETURNING id, name, description, picture, price, category_id, total_stock, available_stock, reserved_stock, low_stock_threshold, allow_negative, created_at, updated_at, version
`

type UpdateProductParams struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Picture     string  `json:"picture"`
	Price       float32 `json:"price"`
	CategoryID  []int32 `json:"categoryID"`
	TotalStock  int32   `json:"totalStock"`
	ID          int32   `json:"id"`
}

//UpdateProduct
//
//  UPDATE products.products
//  SET name = $1, description = $2, picture = $3, price = $4, category_Id = $5, total_stock = $6
//  WHERE id = $7
//  RETURNING id, name, description, picture, price, category_id, total_stock, available_stock, reserved_stock, low_stock_threshold, allow_negative, created_at, updated_at, version
func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) (ProductsProducts, error) {
	row := q.db.QueryRow(ctx, UpdateProduct,
		arg.Name,
		arg.Description,
		arg.Picture,
		arg.Price,
		arg.CategoryID,
		arg.TotalStock,
		arg.ID,
	)
	var i ProductsProducts
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Picture,
		&i.Price,
		&i.CategoryID,
		&i.TotalStock,
		&i.AvailableStock,
		&i.ReservedStock,
		&i.LowStockThreshold,
		&i.AllowNegative,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Version,
	)
	return i, err
}

const UpdateProductsReservedStock = `-- name: UpdateProductsReservedStock :one
UPDATE products.products
SET reserved_stock = reserved_stock + 2
WHERE id = 1
RETURNING id, name, description, picture, price, category_id, total_stock, available_stock, reserved_stock, low_stock_threshold, allow_negative, created_at, updated_at, version
`

// 预留库存（下单时）
//
//  UPDATE products.products
//  SET reserved_stock = reserved_stock + 2
//  WHERE id = 1
//  RETURNING id, name, description, picture, price, category_id, total_stock, available_stock, reserved_stock, low_stock_threshold, allow_negative, created_at, updated_at, version
func (q *Queries) UpdateProductsReservedStock(ctx context.Context) (ProductsProducts, error) {
	row := q.db.QueryRow(ctx, UpdateProductsReservedStock)
	var i ProductsProducts
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Picture,
		&i.Price,
		&i.CategoryID,
		&i.TotalStock,
		&i.AvailableStock,
		&i.ReservedStock,
		&i.LowStockThreshold,
		&i.AllowNegative,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Version,
	)
	return i, err
}
