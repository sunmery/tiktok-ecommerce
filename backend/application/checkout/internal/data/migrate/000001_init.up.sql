CREATE SCHEMA checkout;
SET search_path TO checkout;

-- 支付事务表（记录支付信息，与订单 1:1 关联）
CREATE TABLE checkout.transactions
(
    transaction_id    SERIAL PRIMARY KEY,                     -- 支付事务唯一标识（可来自第三方支付网关）
    order_id          BIGINT,
    credit_card_type  TEXT        NOT NULL,                   -- 卡类型（e.g., Visa/MasterCard）
    brand             TEXT        NOT NULL,                   -- 卡品牌
    country           TEXT        NOT NULL,                   -- 卡所属国家
    credit_card_last4 CHAR(4)     NOT NULL,                   -- 卡号后四位（避免存储敏感信息）
    status            TEXT        NOT NULL DEFAULT 'pending', -- 支付状态（succeeded/failed/pending）
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transactions_order ON checkout.transactions (order_id); -- 按订单查支付
