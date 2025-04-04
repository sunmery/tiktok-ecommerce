CREATE SCHEMA IF NOT EXISTS payments;
SET search_path TO payments;

CREATE TABLE payments.payments
(
    payment_id    UUID PRIMARY KEY,
    order_id      UUID                     NOT NULL, -- 冗余存储，避免跨服务查询
    amount        NUMERIC(12, 2)           NOT NULL,
    currency      VARCHAR(3)               NOT NULL,
    method        VARCHAR(20)              NOT NULL,
    status        VARCHAR(20)              NOT NULL DEFAULT 'PENDING',
    gateway_tx_id VARCHAR(255),                      -- 支付网关交易ID
    metadata      JSONB,                             -- 存储支付扩展信息
    created_at    timestamptz DEFAULT now() NOT NULL,
    updated_at    timestamptz DEFAULT now() NOT NULL
);
