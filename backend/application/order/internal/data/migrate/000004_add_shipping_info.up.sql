-- 主订单表添加物流信息字段
ALTER TABLE orders.sub_orders
    ADD COLUMN tracking_number VARCHAR(100),
    ADD COLUMN carrier VARCHAR(100);