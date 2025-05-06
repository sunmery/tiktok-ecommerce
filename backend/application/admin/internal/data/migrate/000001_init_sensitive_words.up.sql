CREATE SCHEMA IF NOT EXISTS admin;

CREATE TABLE admin.sensitive_words
(
    id         SERIAL PRIMARY KEY,
    created_by UUID                      NOT NULL, -- 添加者
    category   VARCHAR(50)               NOT NULL, -- 类别: 例如, 政治, 色情, 广告
    word       VARCHAR(255) UNIQUE       NOT NULL, -- 敏感词
    level      INT         DEFAULT 1     NOT NULL, -- 敏感级别 用于控制是否不显示该敏感词或警告性言论
    is_active  BOOL        DEFAULT TRUE  NOT NULL,
    created_at timestamptz DEFAULT NOW() NOT NULL,
    updated_at timestamptz DEFAULT NOW() NOT NULL
);

CREATE EXTENSION IF NOT EXISTS pg_trgm; -- 扩展提供了三元模型（trigram）功能，可以极大地加速子串匹配 (LIKE, ILIKE) 和相似度搜索。
-- 创建 GIN 索引
CREATE INDEX IF NOT EXISTS idx_sensitive_words_trgm_gin ON admin.sensitive_words USING gin (word gin_trgm_ops);
