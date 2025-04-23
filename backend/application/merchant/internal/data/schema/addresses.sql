CREATE SCHEMA IF NOT EXISTS merchant;
SET search_path TO merchant;

-- 商家地址表设计
-- CREATE TYPE merchant.address_type AS ENUM ('WAREHOUSE', 'RETURN', 'STORE', 'BILLING');
CREATE TABLE merchant.addresses
(
    id             BIGINT PRIMARY KEY,
    merchant_id    UUID         NOT NULL,
    address_type   VARCHAR(15)  NOT NULL DEFAULT 'WAREHOUSE'
        CHECK ( address_type IN ('WAREHOUSE', 'RETURN', 'STORE', 'BILLING')), -- 地址类型扩展
    contact_person VARCHAR(255) NOT NULL,                                     -- 联系人姓名
    contact_phone  VARCHAR(20)  NOT NULL,                                     -- 联系电话
    street_address TEXT         NOT NULL,
    city           VARCHAR(255) NOT NULL,
    state          VARCHAR(20)  NOT NULL,
    country        VARCHAR(100) NOT NULL,
    zip_code       VARCHAR(20)  NOT NULL,
    is_default     BOOLEAN      NOT NULL DEFAULT FALSE,
    remarks        TEXT         NOT NULL,                                     -- 备注
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
