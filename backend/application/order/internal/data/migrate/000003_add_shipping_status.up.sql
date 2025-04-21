-- 子订单表添加货运状态字段
ALTER TABLE orders.sub_orders
    ADD COLUMN shipping_status VARCHAR(20) NOT NULL DEFAULT 'PENDING_SHIPMENT'
        CHECK (shipping_status IN ('PENDING_SHIPMENT', 'SHIPPED', 'IN_TRANSIT', 'DELIVERED', 'CONFIRMED', 'CANCELLED'));