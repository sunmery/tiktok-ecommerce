CREATE SCHEMA IF NOT EXISTS payments;
SET search_path TO payments;

CREATE FUNCTION uuidv7_sub_ms() RETURNS uuid
AS $$
select encode(
               substring(int8send(floor(t_ms)::int8) from 3) ||
               int2send((7<<12)::int2 | ((t_ms-floor(t_ms))*4096)::int2) ||
               substring(uuid_send(gen_random_uuid()) from 9 for 8)
           , 'hex')::uuid
from (select extract(epoch from clock_timestamp())*1000 as t_ms) s
$$ LANGUAGE sql volatile;

CREATE TABLE payments.payments
(
    payment_id    UUID PRIMARY KEY DEFAULT uuidv7_sub_ms(),
    order_id      UUID                      NOT NULL, -- 冗余存储，避免跨服务查询
    amount        NUMERIC(12, 2)            NOT NULL,
    currency      VARCHAR(3)                NOT NULL,
    method        VARCHAR(20)               NOT NULL,
    status        VARCHAR(20)               NOT NULL DEFAULT 'PENDING',
    gateway_tx_id VARCHAR(255),                       -- 支付网关交易ID
    metadata      JSONB,                              -- 存储支付扩展信息
    created_at    timestamptz DEFAULT now() NOT NULL,
    updated_at    timestamptz DEFAULT now() NOT NULL
);

CREATE INDEX idx_payments_order ON payments.payments (order_id);
CREATE INDEX idx_payments_status ON payments.payments (status);
