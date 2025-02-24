CREATE SCHEMA IF NOT EXISTS categories;
SET search_path TO categories;

CREATE EXTENSION IF NOT EXISTS ltree WITH SCHEMA public;
-- 检查是否安装成功
select *
from pg_extension
where extname = 'ltree';

-- 核心分类表（树形结构核心）
CREATE TABLE IF NOT EXISTS categories.categories
(
    id         BIGINT      NOT NULL,               -- 分类ID（自增序列）
    parent_id  BIGINT      NULL,                   -- 父分类ID
    level      SMALLINT    NOT NULL DEFAULT 1
        CHECK (level BETWEEN 1 AND 3),             -- 层级深度（限制三级）
    path       public.LTREE       NOT NULL,               -- 层级路径（使用PostgreSQL专用ltree类型）
    name       VARCHAR(50) NOT NULL,               -- 分类名称
    sort_order SMALLINT    NOT NULL DEFAULT 0      -- 同级排序（0-32767）
        CHECK (sort_order >= 0),
    is_leaf    BOOLEAN     NOT NULL DEFAULT FALSE, -- 是否为叶子节点
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 创建GIST索引加速ltree查询
CREATE INDEX idx_categories_path_gist ON categories.categories USING GIST (path);
-- B树索引优化常用查询
CREATE INDEX idx_categories_parent ON categories.categories (parent_id);
CREATE INDEX idx_categories_leaf ON categories.categories (is_leaf) WHERE is_leaf = TRUE;

COMMENT ON TABLE categories.categories IS '商品分类主表（ltree路径+闭包表双重优化）';

-- 闭包关系表
CREATE TABLE IF NOT EXISTS categories.category_closure
(
    ancestor   BIGINT   NOT NULL, -- 祖先节点ID
    descendant BIGINT   NOT NULL, -- 后代节点ID
    depth      SMALLINT NOT NULL  -- 层级间隔
        CHECK (depth >= 0),
    PRIMARY KEY (ancestor, descendant)
);

CREATE INDEX idx_closure_descendant ON categories.category_closure (descendant);
COMMENT ON TABLE categories.category_closure IS '分类闭包关系表（存储所有层级关系）';
