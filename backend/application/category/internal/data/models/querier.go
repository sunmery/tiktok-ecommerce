// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package models

import (
	"context"
)

type Querier interface {
	//BatchGetCategories
	//
	//  SELECT
	//      id,
	//      COALESCE(parent_id, 0) AS parent_id,
	//      level,
	//      path::text AS path,
	//      name,
	//      sort_order,
	//      is_leaf,
	//      created_at,
	//      updated_at
	//  FROM categories.categories
	//  WHERE id = ANY($1::bigint[])
	BatchGetCategories(ctx context.Context, ids []int64) ([]BatchGetCategoriesRow, error)
	//CreateCategory
	//
	//  /*
	//  参数说明：
	//  - parent_id: 父分类ID（0表示根）
	//  - name: 分类名称
	//  - sort_order: 排序序号
	//
	//  操作流程：
	//  1. 检查父分类是否存在
	//  2. 插入新分类并生成ltree路径
	//  3. 维护闭包表关系
	//  */
	//
	//  WITH parent_cte AS (
	//      SELECT id, level, path
	//      FROM categories.categories
	//      WHERE id = CASE
	//          WHEN $1::bigint = 0 THEN 1
	//          WHEN $1::bigint IS NULL THEN 1
	//          ELSE $1::bigint
	//      END
	//  ),
	//  valid_parent AS (
	//      SELECT id, level, path
	//      FROM parent_cte
	//      WHERE level < 6
	//      AND EXISTS(SELECT 1 FROM categories.categories WHERE id = (SELECT id FROM parent_cte))
	//  ),
	//  new_category AS (
	//      INSERT INTO categories.categories (
	//          id, parent_id, name, sort_order, level, path, is_leaf
	//      )
	//      SELECT
	//          nextval('categories.categories_id_seq'),
	//          p.id,
	//          $2::varchar(50),
	//          $3::smallint,
	//          p.level + 1,
	//          p.path || currval('categories.categories_id_seq')::text::ltree,
	//          true
	//      FROM valid_parent p
	//      WHERE NOT EXISTS (
	//          SELECT 1 FROM categories.categories
	//          WHERE parent_id = p.id AND name = $2::varchar(50)
	//      )
	//      RETURNING id, parent_id, level, path, name, sort_order, is_leaf, created_at, updated_at
	//  ),
	//  update_leaf AS (
	//      UPDATE categories.categories
	//      SET is_leaf = false
	//      WHERE id = (SELECT parent_id FROM new_category)
	//      AND is_leaf = true
	//  ),
	//  update_current_leaf AS (
	//      UPDATE categories.categories
	//      SET is_leaf = NOT EXISTS (
	//          SELECT 1 FROM categories.categories
	//          WHERE parent_id = (SELECT id FROM new_category)
	//      )
	//      WHERE id = (SELECT id FROM new_category)
	//  ),
	//  closure_insert AS (
	//      INSERT INTO categories.category_closure (ancestor, descendant, depth)
	//      SELECT
	//          c.ancestor,
	//          n.id,
	//          c.depth + 1
	//      FROM categories.category_closure c
	//      CROSS JOIN new_category n
	//      WHERE c.descendant = (SELECT id FROM valid_parent)
	//      UNION ALL
	//      SELECT
	//          n.id,
	//          n.id,
	//          0
	//      FROM new_category n
	//  )
	//  SELECT id, parent_id, level, path, name, sort_order, is_leaf, created_at, updated_at FROM new_category
	CreateCategory(ctx context.Context, arg CreateCategoryParams) (CreateCategoryRow, error)
	//DeleteCategory
	//
	//  WITH deleted_nodes AS (
	//      DELETE FROM categories.categories
	//          WHERE path <@ (SELECT path FROM categories.categories WHERE id = $1)
	//              OR path ~ ($2 || '.*{1,}')::lquery
	//          RETURNING id, parent_id
	//  ),
	//       delete_closure AS (
	//           DELETE FROM categories.category_closure
	//               WHERE descendant IN (SELECT id FROM deleted_nodes)
	//                   OR ancestor IN (SELECT id FROM deleted_nodes)
	//       )
	//  UPDATE categories.categories
	//  SET is_leaf = (
	//      SELECT NOT EXISTS (
	//          SELECT 1 FROM categories.categories
	//          WHERE parent_id = (SELECT parent_id FROM deleted_nodes LIMIT 1)
	//            AND id != (SELECT id FROM deleted_nodes LIMIT 1)
	//      )
	//  )
	//  WHERE id = (SELECT parent_id FROM deleted_nodes LIMIT 1)
	//    AND (SELECT parent_id FROM deleted_nodes LIMIT 1) IS NOT NULL
	DeleteCategory(ctx context.Context, arg DeleteCategoryParams) error
	//GetCategory
	//
	//  SELECT
	//      id,
	//      COALESCE(parent_id, 0) AS parent_id,  -- 将NULL转换为0返回给proto
	//      level,
	//      path::text AS path,
	//      name,
	//      sort_order,
	//      is_leaf,
	//      created_at,
	//      updated_at
	//  FROM categories.categories WHERE id = $1
	GetCategory(ctx context.Context, id int64) (GetCategoryRow, error)
	//GetCategoryPath
	//
	//  SELECT
	//      c.id,
	//      COALESCE(c.parent_id, 0) AS parent_id,
	//      c.level,
	//      c.path::text,
	//      c.name,
	//      c.sort_order,
	//      c.is_leaf,
	//      c.created_at,
	//      c.updated_at
	//  FROM categories.category_closure cc
	//           JOIN categories.categories c ON cc.ancestor = c.id
	//  WHERE cc.descendant = $1
	//  ORDER BY cc.depth DESC
	GetCategoryPath(ctx context.Context, descendant int64) ([]GetCategoryPathRow, error)
	//GetClosureRelations
	//
	//  SELECT ancestor, descendant, depth FROM categories.category_closure
	//  WHERE descendant = $1
	GetClosureRelations(ctx context.Context, descendant int64) ([]CategoriesCategoryClosure, error)
	//GetLeafCategories
	//
	//  SELECT
	//      id,
	//      COALESCE(parent_id, 0) AS parent_id,
	//      level,
	//      path::text,
	//      name,
	//      sort_order,
	//      is_leaf,
	//      created_at,
	//      updated_at
	//  FROM categories.categories
	//  WHERE level != 0
	GetLeafCategories(ctx context.Context) ([]GetLeafCategoriesRow, error)
	//GetSubTree
	//
	//  SELECT
	//      c.id,
	//      COALESCE(c.parent_id, 0) AS parent_id,
	//      c.level,
	//      c.path::text,
	//      c.name,
	//      c.sort_order,
	//      c.is_leaf,
	//      c.created_at,
	//      c.updated_at
	//  FROM categories.category_closure cc
	//           JOIN categories.categories c ON cc.descendant = c.id
	//  WHERE cc.ancestor = $1 AND cc.depth >= 0
	//  ORDER BY cc.depth
	GetSubTree(ctx context.Context, ancestor int64) ([]GetSubTreeRow, error)
	//UpdateCategory
	//
	//  UPDATE categories.categories
	//  SET name = $2, updated_at = NOW()
	//  WHERE id = $1
	UpdateCategory(ctx context.Context, arg UpdateCategoryParams) error
	//UpdateClosureDepth
	//
	//  UPDATE categories.category_closure
	//  SET depth = depth + $2
	//  WHERE descendant = $1
	UpdateClosureDepth(ctx context.Context, arg UpdateClosureDepthParams) error
}

var _ Querier = (*Queries)(nil)
