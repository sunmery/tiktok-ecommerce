-- 主订单表
ALTER TABLE orders.orders
    ADD COLUMN payment_status VARCHAR(20) NOT NULL DEFAULT 'PENDING'
        CHECK (payment_status IN ('PENDING', 'PAID', 'CANCELLED', 'FAILED', 'CANCELLED'));

-- 子订单表
ALTER TABLE orders.sub_orders
    ADD COLUMN payment_status VARCHAR(20) NOT NULL DEFAULT 'PENDING';
