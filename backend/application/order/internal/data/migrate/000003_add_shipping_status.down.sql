-- 删除主订单表中的货运状态字段
ALTER TABLE orders.orders
    DROP COLUMN shipping_status;