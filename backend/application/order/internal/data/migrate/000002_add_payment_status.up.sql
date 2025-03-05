-- 主订单表
ALTER TABLE orders.orders
    ADD COLUMN payment_status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (payment_status IN ('pending', 'paid', 'cancelled', 'failed', 'cancelled'));

-- 子订单表
ALTER TABLE orders.sub_orders
    ADD COLUMN payment_status VARCHAR(20) NOT NULL DEFAULT 'pending';
