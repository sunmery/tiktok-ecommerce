CREATE SCHEMA IF NOT EXISTS admin;

CREATE TABLE admin.sensitive_words
(
    id         SERIAL PRIMARY KEY,
    created_by UUID                      NOT NULL, -- 添加者
    category   VARCHAR(50)               NOT NULL, -- 类别
    word       VARCHAR(255) UNIQUE       NOT NULL, -- 敏感词
    level      INT         DEFAULT 1     NOT NULL, -- 敏感级别
    is_active  BOOL        DEFAULT TRUE  NOT NULL, -- 是否激活
    created_at timestamptz DEFAULT NOW() NOT NULL, -- 创建时间
    updated_at timestamptz DEFAULT NOW() NOT NULL  -- 更新时间
);
COMMENT ON TABLE admin.sensitive_words IS '敏感词表';

CREATE EXTENSION IF NOT EXISTS pg_trgm;
-- 扩展提供了三元模型（trigram）功能，可以极大地加速子串匹配 (LIKE, ILIKE) 和相似度搜索。
-- 创建 GIN 索引
CREATE INDEX IF NOT EXISTS idx_sensitive_words_trgm_gin ON admin.sensitive_words USING gin (word gin_trgm_ops);
