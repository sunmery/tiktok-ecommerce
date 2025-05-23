-- 创建物流表
CREATE TABLE orders.shipping_info
(
    id               BIGINT PRIMARY KEY,
    sub_order_id     BIGINT         NOT NULL,              -- 关联子订单ID
    merchant_id      UUID           NOT NULL,              -- 商家 id
    tracking_number  VARCHAR(100)   NOT NULL,              -- 物流单号
    carrier          VARCHAR(100)   NOT NULL,              -- 承运商
    shipping_status  VARCHAR(15)    NOT NULL DEFAULT 'WAIT_COMMAND'
        CHECK ( shipping_status IN
                ('WAIT_COMMAND', 'PENDING_SHIPMENT', 'SHIPPED', 'IN_TRANSIT', 'DELIVERED', 'CONFIRMED',
                 'CANCELLED')),                            -- 物流状态
    delivery         TIMESTAMPTZ    NULL,                  -- 送达时间
    shipping_address JSONB          NOT NULL DEFAULT '{}', -- 发货地址信息
    receiver_address JSONB          NOT NULL DEFAULT '{}', -- 收货地址信息
    shipping_fee     NUMERIC(10, 2) NOT NULL DEFAULT 0.00, -- 运费
    created_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

-- 添加索引
CREATE INDEX idx_shipping_sub_order_id ON orders.shipping_info (sub_order_id);
CREATE INDEX idx_shipping_tracking_number ON orders.shipping_info (tracking_number);
