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

-- name: CreateCategory :one
WITH parent_cte AS (
    SELECT id, level, path
    FROM categories.categories
    WHERE id = CASE
        WHEN @parent_id::bigint = 0 THEN 1
        WHEN @parent_id::bigint IS NULL THEN 1
        ELSE @parent_id::bigint
    END
),
valid_parent AS (
    SELECT *
    FROM parent_cte
    WHERE level < 6
    AND EXISTS(SELECT 1 FROM categories.categories WHERE id = (SELECT id FROM parent_cte))
),
new_category AS (
    INSERT INTO categories.categories (
        id, parent_id, name, sort_order, level, path, is_leaf
    )
    SELECT
        nextval('categories.categories_id_seq'),
        p.id,
        @name::varchar(50),
        @sort_order::smallint,
        p.level + 1,
        p.path || currval('categories.categories_id_seq')::text::ltree,
        true
    FROM valid_parent p
    WHERE NOT EXISTS (
        SELECT 1 FROM categories.categories
        WHERE parent_id = p.id AND name = @name::varchar(50)
    )
    RETURNING *
),
update_leaf AS (
    UPDATE categories.categories
    SET is_leaf = false
    WHERE id = (SELECT parent_id FROM new_category)
    AND is_leaf = true
),
update_current_leaf AS (
    UPDATE categories.categories
    SET is_leaf = NOT EXISTS (
        SELECT 1 FROM categories.categories
        WHERE parent_id = (SELECT id FROM new_category)
    )
    WHERE id = (SELECT id FROM new_category)
),
closure_insert AS (
    INSERT INTO categories.category_closure (ancestor, descendant, depth)
    SELECT
        c.ancestor,
        n.id,
        c.depth + 1
    FROM categories.category_closure c
    CROSS JOIN new_category n
    WHERE c.descendant = (SELECT id FROM valid_parent)
    UNION ALL
    SELECT
        n.id,
        n.id,
        0
    FROM new_category n
)
SELECT * FROM new_category;

-- name: GetCategory :one
SELECT
    id,
    COALESCE(parent_id, 0) AS parent_id,  -- 将NULL转换为0返回给proto
    level,
    path::text AS path,
    name,
    sort_order,
    is_leaf,
    created_at,
    updated_at
FROM categories.categories WHERE id = $1;

-- name: UpdateCategory :exec
UPDATE categories.categories
SET name = $2, updated_at = NOW()
WHERE id = $1;

-- name: DeleteCategory :exec
WITH deleted_nodes AS (
    DELETE FROM categories.categories
        WHERE path <@ (SELECT path FROM categories.categories WHERE id = @id)
            OR path ~ (@path || '.*{1,}')::lquery
        RETURNING id, parent_id
),
     delete_closure AS (
         DELETE FROM categories.category_closure
             WHERE descendant IN (SELECT id FROM deleted_nodes)
                 OR ancestor IN (SELECT id FROM deleted_nodes)
     )
UPDATE categories.categories
SET is_leaf = (
    SELECT NOT EXISTS (
        SELECT 1 FROM categories.categories
        WHERE parent_id = (SELECT parent_id FROM deleted_nodes LIMIT 1)
          AND id != (SELECT id FROM deleted_nodes LIMIT 1)
    )
)
WHERE id = (SELECT parent_id FROM deleted_nodes LIMIT 1)
  AND (SELECT parent_id FROM deleted_nodes LIMIT 1) IS NOT NULL;

-- name: GetSubTree :many
WITH RECURSIVE category_tree AS (
    -- 基本情况：直接获取所有子节点作为起点
    SELECT
        c.id,
        c.parent_id,
        c.level,
        c.path,
        c.name,
        c.sort_order,
        c.is_leaf,
        c.created_at,
        c.updated_at
    FROM categories.categories c
    WHERE c.parent_id = $1
    
    UNION ALL
    
    -- 递归情况：获取所有直接子节点
    SELECT
        c.id,
        c.parent_id,
        c.level,
        c.path,
        c.name,
        c.sort_order,
        c.is_leaf,
        c.created_at,
        c.updated_at
    FROM categories.categories c
    JOIN category_tree ct ON c.parent_id = ct.id
)
SELECT
    id,
    COALESCE(parent_id, 0) AS parent_id,
    level,
    path::text AS cpath,
    name,
    sort_order,
    is_leaf,
    created_at,
    updated_at
FROM category_tree
ORDER BY level, id;

-- name: GetDirectSubCategories :many
SELECT
    id,
    COALESCE(parent_id, 0) AS parent_id,
    level,
    path::text AS cpath,
    name,
    sort_order,
    is_leaf,
    created_at,
    updated_at
FROM categories.categories
WHERE parent_id = $1
ORDER BY sort_order, id;

-- name: GetCategoryPath :many
SELECT
    c.id,
    COALESCE(c.parent_id, 0) AS parent_id,
    c.level,
    c.path::text,
    c.name,
    c.sort_order,
    c.is_leaf,
    c.created_at,
    c.updated_at
FROM categories.category_closure cc
         JOIN categories.categories c ON cc.ancestor = c.id
WHERE cc.descendant = $1
ORDER BY cc.depth DESC;

-- name: GetLeafCategories :many
SELECT
    id,
    COALESCE(parent_id, 0) AS parent_id,
    level,
    path::text,
    name,
    sort_order,
    is_leaf,
    created_at,
    updated_at
FROM categories.categories
WHERE level != 0;

-- name: GetClosureRelations :many
SELECT * FROM categories.category_closure
WHERE descendant = $1;

-- name: UpdateClosureDepth :exec
UPDATE categories.category_closure
SET depth = depth + $2
WHERE descendant = $1;

-- name: BatchGetCategories :many
SELECT
    id,
    COALESCE(parent_id, 0) AS parent_id,
    level,
    path::text AS path,
    name,
    sort_order,
    is_leaf,
    created_at,
    updated_at
FROM categories.categories
WHERE id = ANY(@ids::bigint[]);