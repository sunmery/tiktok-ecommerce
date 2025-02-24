CREATE SCHEMA IF NOT EXISTS categories;
SET search_path TO categories;
-- 使用ltree类型
/*
优势：
1. 原生支持路径查询：'0.123.456' <@ '0.123'
2. 提供多种树操作符：@>（包含）、<@（被包含）
3. 支持GIST索引加速树查询
*/

-- 移除外键后的数据完整性保障方案
/*
在Go程序中需要：
1. 插入时验证parent_id存在性
2. 删除时维护闭包表一致性
3. 使用事务保证操作原子性
*/

-- 在public共享, 其他的 schema 也可以使用
-- DB_SOURCE使用: search_path=categories,public
CREATE EXTENSION IF NOT EXISTS ltree WITH SCHEMA public;
-- 检查是否安装成功
select *
from pg_extension
where extname = 'ltree';

-- 核心分类表（树形结构核心）
CREATE TABLE IF NOT EXISTS categories.categories
(
    id         BIGINT,                             -- 分类ID（自增序列）
    parent_id  BIGINT      NULL,                   -- 父分类ID
    level      SMALLINT    NOT NULL DEFAULT 1
        CHECK (level BETWEEN 1 AND 3),             -- 层级深度（限制三级）
    path       LTREE       NOT NULL,               -- 层级路径（使用PostgreSQL专用ltree类型）
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
-- 修改分类表的 level 约束
-- ALTER TABLE categories.categories
-- DROP CONSTRAINT level_check;

ALTER TABLE categories.categories
    ADD CONSTRAINT level_check CHECK (level BETWEEN 1 AND 4);

-- 修改闭包表的 depth 约束
-- ALTER TABLE categories.category_closure
-- DROP CONSTRAINT depth_check;

ALTER TABLE categories.category_closure
    ADD CONSTRAINT depth_check CHECK (depth BETWEEN 0 AND 3); -- 0 到 3 表示四层
