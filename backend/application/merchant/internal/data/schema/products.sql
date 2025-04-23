CREATE TABLE products.products
(
    id               UUID                  DEFAULT uuidv7_sub_ms(),
    merchant_id      UUID         NOT NULL, -- 分片键（必须）
    name             VARCHAR(255) NOT NULL,
    description      TEXT,
    price            NUMERIC(15, 2) CHECK (price >= 0),
    status           SMALLINT     NOT NULL DEFAULT 1,
    current_audit_id UUID,
    category_id      int8         NOT NULL, -- 商品分类 ID
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ,
    PRIMARY KEY (merchant_id, id)           -- Citus 分片表需要包含分片键在PK中
);
