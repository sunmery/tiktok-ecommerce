CREATE SCHEMA IF NOT EXISTS payments;
SET search_path TO payments;

CREATE TABLE payments.payments
(
    id          BIGINT PRIMARY KEY,
    order_id    BIGINT                    NOT NULL,                   -- 订单ID, 冗余存储，避免跨服务查询
    user_id     UUID                      NOT NULL,                   -- 用户ID，便于后续操作
    amount      NUMERIC(12, 2)            NOT NULL,                   -- 订单金额
    currency    VARCHAR(3)                NOT NULL,                   -- 货币类型, 例如: CNY, USD
    method      VARCHAR(20)               NOT NULL,                   -- 支付方式, 例如: 支付宝, 微信 等等
    status      VARCHAR(20)               NOT NULL DEFAULT 'PENDING', -- 支付状态
    subject     VARCHAR(255)              NOT NULL,                   -- 支付标题
    trade_no    VARCHAR(255)              NOT NULL,                   -- 商户订单号
    metadata    JSONB                     NULL,                       -- 存储支付扩展信息
    created_at  timestamptz DEFAULT now() NOT NULL,
    updated_at  timestamptz DEFAULT now() NOT NULL
);