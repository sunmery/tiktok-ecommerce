CREATE SCHEMA IF NOT EXISTS merchant;
SET search_path TO merchant;

-- 商家地址表设计
CREATE TABLE merchant.addresses
(
    id             BIGINT PRIMARY KEY,
    merchant_id    UUID         NOT NULL,
    address_type   VARCHAR(20)  NOT NULL DEFAULT 'WAREHOUSE' CHECK (
        address_type IN ('WAREHOUSE', 'RETURN', 'STORE', 'BILLING')
        ),                                -- 地址类型扩展
    contact_person VARCHAR(255) NOT NULL, -- 联系人姓名
    contact_phone  VARCHAR(20)  NOT NULL, -- 联系电话
    street_address TEXT         NOT NULL,
    city           VARCHAR(255) NOT NULL,
    state          VARCHAR(20)  NOT NULL,
    country        VARCHAR(100) NOT NULL,
    zip_code       VARCHAR(20)  NOT NULL,
    is_default     BOOLEAN      NOT NULL DEFAULT FALSE,
    remarks        TEXT         NOT NULL, -- 备注
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- 索引优化方案
CREATE INDEX idx_merchant_addr_main ON merchant.addresses (merchant_id)
    WHERE is_default = true; -- 快速查询默认地址

CREATE INDEX idx_merchant_addr_type ON merchant.addresses
    USING BRIN (merchant_id, address_type); -- 时间范围查询优化
