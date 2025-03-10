// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: category.sql

package models

import (
	"context"
)

const CreateCategory = `-- name: CreateCategory :one
/*
参数说明：
- parent_id: 父分类ID（0表示根）
- name: 分类名称
- sort_order: 排序序号

操作流程：
1. 检查父分类是否存在
2. 插入新分类并生成ltree路径
3. 维护闭包表关系
*/

WITH root_check AS (
  INSERT INTO categories.categories (id, parent_id, level, path, name, sort_order, is_leaf)
  VALUES (0, 0, 1, 'root'::public.ltree, 'Root', 0, FALSE)
  ON CONFLICT (id) DO NOTHING
),
parent_info AS (
  SELECT
    COALESCE(c.id, 0) AS effective_parent_id,
    COALESCE(c.path, 'root'::public.ltree) AS parent_path,
    COALESCE(c.level, 0) AS parent_level
  FROM (SELECT $1::BIGINT AS pid) AS input
  LEFT JOIN categories.categories c ON c.id = input.pid
),
level_validation AS (
  SELECT
    effective_parent_id,
    parent_path,
    CASE
      WHEN parent_level >= 4 THEN NULL -- 父节点已经是4层，不允许新增子节点
      ELSE parent_level + 1
    END AS new_level
  FROM parent_info
),
insert_main AS (
  INSERT INTO categories.categories (
    parent_id, level, path, name, sort_order, is_leaf
  ) SELECT
    lv.effective_parent_id,
    lv.new_level,
    CASE
      WHEN lv.parent_path OPERATOR(public.=) 'root'::public.ltree
      THEN lv.parent_path || ('node_' || gen_random_uuid())::public.ltree
      ELSE lv.parent_path || (REPLACE(gen_random_uuid()::text, '-', '_'))::public.ltree
    END,
    $2,   -- Name 参数
    $3,   -- SortOrder 参数
    CASE WHEN lv.new_level = 4 THEN TRUE ELSE FALSE END -- 第四层为叶子节点
  FROM level_validation lv
  WHERE lv.new_level IS NOT NULL
  RETURNING id, parent_id, level, path, name, sort_order, is_leaf, created_at, updated_at
),
update_parent_leaf AS (
  UPDATE categories.categories
  SET is_leaf = FALSE
  WHERE id = (SELECT effective_parent_id FROM parent_info)
  AND is_leaf = TRUE
)
INSERT INTO categories.category_closure (ancestor, descendant, depth)
SELECT
  cc.ancestor,
  im.id,
  cc.depth + 1
FROM insert_main im
JOIN categories.category_closure cc ON cc.descendant = im.parent_id
UNION ALL
SELECT
  im.id,
  im.id,
  0
FROM insert_main im
RETURNING descendant
`

type CreateCategoryParams struct {
	ParentID  int64  `json:"parent_id"`
	Name      string `json:"name"`
	SortOrder int16  `json:"sort_order"`
}

// CreateCategory
//
//	/*
//	参数说明：
//	- parent_id: 父分类ID（0表示根）
//	- name: 分类名称
//	- sort_order: 排序序号
//
//	操作流程：
//	1. 检查父分类是否存在
//	2. 插入新分类并生成ltree路径
//	3. 维护闭包表关系
//	*/
//
//	WITH root_check AS (
//	  INSERT INTO categories.categories (id, parent_id, level, path, name, sort_order, is_leaf)
//	  VALUES (0, 0, 1, 'root'::public.ltree, 'Root', 0, FALSE)
//	  ON CONFLICT (id) DO NOTHING
//	),
//	parent_info AS (
//	  SELECT
//	    COALESCE(c.id, 0) AS effective_parent_id,
//	    COALESCE(c.path, 'root'::public.ltree) AS parent_path,
//	    COALESCE(c.level, 0) AS parent_level
//	  FROM (SELECT $1::BIGINT AS pid) AS input
//	  LEFT JOIN categories.categories c ON c.id = input.pid
//	),
//	level_validation AS (
//	  SELECT
//	    effective_parent_id,
//	    parent_path,
//	    CASE
//	      WHEN parent_level >= 4 THEN NULL -- 父节点已经是4层，不允许新增子节点
//	      ELSE parent_level + 1
//	    END AS new_level
//	  FROM parent_info
//	),
//	insert_main AS (
//	  INSERT INTO categories.categories (
//	    parent_id, level, path, name, sort_order, is_leaf
//	  ) SELECT
//	    lv.effective_parent_id,
//	    lv.new_level,
//	    CASE
//	      WHEN lv.parent_path OPERATOR(public.=) 'root'::public.ltree
//	      THEN lv.parent_path || ('node_' || gen_random_uuid())::public.ltree
//	      ELSE lv.parent_path || (REPLACE(gen_random_uuid()::text, '-', '_'))::public.ltree
//	    END,
//	    $2,   -- Name 参数
//	    $3,   -- SortOrder 参数
//	    CASE WHEN lv.new_level = 4 THEN TRUE ELSE FALSE END -- 第四层为叶子节点
//	  FROM level_validation lv
//	  WHERE lv.new_level IS NOT NULL
//	  RETURNING id, parent_id, level, path, name, sort_order, is_leaf, created_at, updated_at
//	),
//	update_parent_leaf AS (
//	  UPDATE categories.categories
//	  SET is_leaf = FALSE
//	  WHERE id = (SELECT effective_parent_id FROM parent_info)
//	  AND is_leaf = TRUE
//	)
//	INSERT INTO categories.category_closure (ancestor, descendant, depth)
//	SELECT
//	  cc.ancestor,
//	  im.id,
//	  cc.depth + 1
//	FROM insert_main im
//	JOIN categories.category_closure cc ON cc.descendant = im.parent_id
//	UNION ALL
//	SELECT
//	  im.id,
//	  im.id,
//	  0
//	FROM insert_main im
//	RETURNING descendant
func (q *Queries) CreateCategory(ctx context.Context, arg CreateCategoryParams) (int64, error) {
	row := q.db.QueryRow(ctx, CreateCategory, arg.ParentID, arg.Name, arg.SortOrder)
	var descendant int64
	err := row.Scan(&descendant)
	return descendant, err
}

const DeleteCategory = `-- name: DeleteCategory :exec
WITH deleted AS (
    DELETE FROM categories.categories
        WHERE id = $1
        RETURNING path
)
DELETE FROM categories.category_closure
WHERE descendant IN (
    SELECT descendant
    FROM categories.category_closure
    WHERE ancestor = $1
)
`

// DeleteCategory
//
//	WITH deleted AS (
//	    DELETE FROM categories.categories
//	        WHERE id = $1
//	        RETURNING path
//	)
//	DELETE FROM categories.category_closure
//	WHERE descendant IN (
//	    SELECT descendant
//	    FROM categories.category_closure
//	    WHERE ancestor = $1
//	)
func (q *Queries) DeleteCategory(ctx context.Context, id *int64) error {
	_, err := q.db.Exec(ctx, DeleteCategory, id)
	return err
}

const DeleteClosureRelations = `-- name: DeleteClosureRelations :exec

DELETE FROM categories.category_closure
WHERE descendant IN (
    SELECT descendant
    FROM categories.category_closure
    WHERE ancestor = $1
)
`

// 确保深度不超过 3
// 删除指定分类及其所有后代节点的闭包关系
//
//	DELETE FROM categories.category_closure
//	WHERE descendant IN (
//	    SELECT descendant
//	    FROM categories.category_closure
//	    WHERE ancestor = $1
//	)
func (q *Queries) DeleteClosureRelations(ctx context.Context, categoryID *int64) error {
	_, err := q.db.Exec(ctx, DeleteClosureRelations, categoryID)
	return err
}

const GetCategoryByID = `-- name: GetCategoryByID :one
SELECT id, parent_id, level, path, name, sort_order, is_leaf, created_at, updated_at FROM categories.categories
WHERE id = $1 LIMIT 1
`

// GetCategoryByID
//
//	SELECT id, parent_id, level, path, name, sort_order, is_leaf, created_at, updated_at FROM categories.categories
//	WHERE id = $1 LIMIT 1
func (q *Queries) GetCategoryByID(ctx context.Context, id int64) (CategoriesCategories, error) {
	row := q.db.QueryRow(ctx, GetCategoryByID, id)
	var i CategoriesCategories
	err := row.Scan(
		&i.ID,
		&i.ParentID,
		&i.Level,
		&i.Path,
		&i.Name,
		&i.SortOrder,
		&i.IsLeaf,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const GetCategoryPath = `-- name: GetCategoryPath :many
SELECT ancestor.id, ancestor.parent_id, ancestor.level, ancestor.path, ancestor.name, ancestor.sort_order, ancestor.is_leaf, ancestor.created_at, ancestor.updated_at
FROM categories.category_closure cc
         JOIN categories.categories ancestor ON cc.ancestor = ancestor.id
WHERE cc.descendant = $1
ORDER BY cc.depth DESC
`

// GetCategoryPath
//
//	SELECT ancestor.id, ancestor.parent_id, ancestor.level, ancestor.path, ancestor.name, ancestor.sort_order, ancestor.is_leaf, ancestor.created_at, ancestor.updated_at
//	FROM categories.category_closure cc
//	         JOIN categories.categories ancestor ON cc.ancestor = ancestor.id
//	WHERE cc.descendant = $1
//	ORDER BY cc.depth DESC
func (q *Queries) GetCategoryPath(ctx context.Context, categoryID int64) ([]CategoriesCategories, error) {
	rows, err := q.db.Query(ctx, GetCategoryPath, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CategoriesCategories
	for rows.Next() {
		var i CategoriesCategories
		if err := rows.Scan(
			&i.ID,
			&i.ParentID,
			&i.Level,
			&i.Path,
			&i.Name,
			&i.SortOrder,
			&i.IsLeaf,
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

const GetClosureRelations = `-- name: GetClosureRelations :many
SELECT ancestor, descendant, depth FROM categories.category_closure
WHERE descendant = $1
`

// GetClosureRelations
//
//	SELECT ancestor, descendant, depth FROM categories.category_closure
//	WHERE descendant = $1
func (q *Queries) GetClosureRelations(ctx context.Context, categoryID int64) ([]CategoriesCategoryClosure, error) {
	rows, err := q.db.Query(ctx, GetClosureRelations, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CategoriesCategoryClosure
	for rows.Next() {
		var i CategoriesCategoryClosure
		if err := rows.Scan(&i.Ancestor, &i.Descendant, &i.Depth); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetLeafCategories = `-- name: GetLeafCategories :many
SELECT id, parent_id, level, path, name, sort_order, is_leaf, created_at, updated_at FROM categories.categories
WHERE is_leaf = TRUE AND level = 4
`

// GetLeafCategories
//
//	SELECT id, parent_id, level, path, name, sort_order, is_leaf, created_at, updated_at FROM categories.categories
//	WHERE is_leaf = TRUE AND level = 4
func (q *Queries) GetLeafCategories(ctx context.Context) ([]CategoriesCategories, error) {
	rows, err := q.db.Query(ctx, GetLeafCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CategoriesCategories
	for rows.Next() {
		var i CategoriesCategories
		if err := rows.Scan(
			&i.ID,
			&i.ParentID,
			&i.Level,
			&i.Path,
			&i.Name,
			&i.SortOrder,
			&i.IsLeaf,
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

const GetSubTree = `-- name: GetSubTree :many
/*
级联删除策略：
1. 根据闭包表找到所有后代节点
2. 删除所有相关闭包关系
*/


SELECT c.id, c.parent_id, c.level, c.path, c.name, c.sort_order, c.is_leaf, c.created_at, c.updated_at
FROM categories.categories c
WHERE c.path <@ (SELECT path FROM categories.categories WHERE id = $1)
ORDER BY c.path
`

// GetSubTree
//
//	/*
//	级联删除策略：
//	1. 根据闭包表找到所有后代节点
//	2. 删除所有相关闭包关系
//	*/
//
//
//	SELECT c.id, c.parent_id, c.level, c.path, c.name, c.sort_order, c.is_leaf, c.created_at, c.updated_at
//	FROM categories.categories c
//	WHERE c.path <@ (SELECT path FROM categories.categories WHERE id = $1)
//	ORDER BY c.path
func (q *Queries) GetSubTree(ctx context.Context, rootID *int64) ([]CategoriesCategories, error) {
	rows, err := q.db.Query(ctx, GetSubTree, rootID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CategoriesCategories
	for rows.Next() {
		var i CategoriesCategories
		if err := rows.Scan(
			&i.ID,
			&i.ParentID,
			&i.Level,
			&i.Path,
			&i.Name,
			&i.SortOrder,
			&i.IsLeaf,
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

const UpdateCategoryName = `-- name: UpdateCategoryName :exec
UPDATE categories.categories
SET name = $1, updated_at = NOW()
WHERE id = $2
`

type UpdateCategoryNameParams struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

// UpdateCategoryName
//
//	UPDATE categories.categories
//	SET name = $1, updated_at = NOW()
//	WHERE id = $2
func (q *Queries) UpdateCategoryName(ctx context.Context, arg UpdateCategoryNameParams) error {
	_, err := q.db.Exec(ctx, UpdateCategoryName, arg.Name, arg.ID)
	return err
}

const UpdateClosureDepth = `-- name: UpdateClosureDepth :exec
UPDATE categories.category_closure
SET depth = depth + $1
WHERE descendant IN (
    SELECT descendant
    FROM categories.category_closure
    WHERE ancestor = $2
)
AND depth + $1 <= 3
`

type UpdateClosureDepthParams struct {
	Delta      *int16 `json:"delta"`
	CategoryID *int64 `json:"category_id"`
}

// UpdateClosureDepth
//
//	UPDATE categories.category_closure
//	SET depth = depth + $1
//	WHERE descendant IN (
//	    SELECT descendant
//	    FROM categories.category_closure
//	    WHERE ancestor = $2
//	)
//	AND depth + $1 <= 3
func (q *Queries) UpdateClosureDepth(ctx context.Context, arg UpdateClosureDepthParams) error {
	_, err := q.db.Exec(ctx, UpdateClosureDepth, arg.Delta, arg.CategoryID)
	return err
}

const UpdateParentLeafStatus = `-- name: UpdateParentLeafStatus :exec
UPDATE categories.categories
SET
    is_leaf = NOT EXISTS (
        SELECT 1
        FROM categories
        WHERE parent_id = $1
        LIMIT 1
    ),
    updated_at = NOW()
WHERE id = $1
`

// 更新父分类的叶子节点状态
//
//	UPDATE categories.categories
//	SET
//	    is_leaf = NOT EXISTS (
//	        SELECT 1
//	        FROM categories
//	        WHERE parent_id = $1
//	        LIMIT 1
//	    ),
//	    updated_at = NOW()
//	WHERE id = $1
func (q *Queries) UpdateParentLeafStatus(ctx context.Context, parentID *int64) error {
	_, err := q.db.Exec(ctx, UpdateParentLeafStatus, parentID)
	return err
}
