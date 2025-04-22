-- 子订单表添加地址字段
ALTER TABLE orders.sub_orders
    ADD merchant_address JSONB;
