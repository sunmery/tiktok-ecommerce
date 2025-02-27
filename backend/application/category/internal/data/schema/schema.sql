CREATE SCHEMA IF NOT EXISTS categories;
SET search_path TO categories;

CREATE EXTENSION IF NOT EXISTS ltree;

CREATE SEQUENCE categories.categories_id_seq
    START WITH 2  -- 初始值设为2（根节点已用1）
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE categories.categories (
    id         BIGINT PRIMARY KEY DEFAULT nextval('categories_id_seq'),
    parent_id  BIGINT,
    level      SMALLINT NOT NULL CHECK (level BETWEEN 0 AND 3),
    path       LTREE NOT NULL,
    name       VARCHAR(50) NOT NULL,
    sort_order SMALLINT NOT NULL DEFAULT 0 CHECK (sort_order >= 0),
    is_leaf    BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (parent_id, name)
);

CREATE INDEX idx_categories_path_gist ON categories.categories USING GIST (path);

-- 闭包表
CREATE TABLE categories.category_closure (
    ancestor   BIGINT NOT NULL REFERENCES categories.categories(id),
    descendant BIGINT NOT NULL REFERENCES categories.categories(id),
    depth      SMALLINT NOT NULL CHECK (depth >= 0),
    PRIMARY KEY (ancestor, descendant)
);
