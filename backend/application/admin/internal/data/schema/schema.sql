CREATE TABLE admin.sensitive_words
(
    id         SERIAL PRIMARY KEY,
    category   VARCHAR(50)               NOT NULL, -- 类别: 例如, 政治, 色情, 广告
    word       VARCHAR(255) UNIQUE       NOT NULL, -- 敏感词
    level      INT         DEFAULT 1     NOT NULL, -- 敏感级别 用于控制是否不显示该敏感词或警告性言论
    is_active  BOOLEAN     DEFAULT TRUE  NOT NULL,
    created_at timestamptz DEFAULT NOW() NOT NULL,
    updated_at timestamptz DEFAULT NOW() NOT NULL
);