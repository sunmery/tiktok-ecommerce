CREATE SCHEMA IF NOT EXISTS payments;
SET search_path TO payments;

CREATE TABLE payments.payments
(
    id            BIGINT PRIMARY KEY,
    order_id      BIGINT                    NOT NULL, -- 冗余存储，避免跨服务查询
    user_id       UUID                      NOT NULL, -- 用户ID，便于后续操作
    amount        NUMERIC(12, 2)            NOT NULL,
    currency      VARCHAR(3)                NOT NULL,
    method        VARCHAR(20)               NOT NULL,
    status        VARCHAR(20)               NOT NULL DEFAULT 'PENDING',
    subject       VARCHAR(255)              NOT NULL, -- 支付标题
    trade_no      VARCHAR(255)              NOT NULL, -- 商户订单号
    gateway_tx_id VARCHAR(255)              NOT NULL, -- 支付网关交易ID
    pay_url       TEXT                      NOT NULL, -- 支付链接
    metadata      JSONB                     NULL,     -- 存储支付扩展信息
    created_at    timestamptz DEFAULT now() NOT NULL,
    updated_at    timestamptz DEFAULT now() NOT NULL
);

CREATE INDEX idx_payments_order ON payments.payments (order_id);
CREATE INDEX idx_payments_status ON payments.payments (status);
CREATE INDEX idx_payments_user_id ON payments.payments (user_id);
CREATE INDEX idx_payments_trade_no ON payments.payments (trade_no);
