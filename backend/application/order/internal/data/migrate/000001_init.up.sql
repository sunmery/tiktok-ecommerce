CREATE SCHEMA IF NOT EXISTS orders;


-- 主订单表，记录订单全局信息
CREATE TABLE orders.main_orders
(
    order_id     SERIAL PRIMARY KEY,            -- 主订单唯一标识
    user_id      INT            NOT NULL,       -- 用户ID（关联用户表）
    total_amount NUMERIC(10, 2) NOT NULL,-- 订单总金额（含缺货商品）
    status       VARCHAR(20) DEFAULT 'pending', -- 状态：pending/paid/cancelled
    created_at   TIMESTAMP   DEFAULT NOW()      -- 订单创建时间
);

-- 子订单表，按商家拆分订单
CREATE TABLE sub_orders
(
    sub_order_id  SERIAL PRIMARY KEY,                           -- 子订单唯一标识
    main_order_id INT REFERENCES orders.main_orders (order_id), -- 关联主订单
    seller_id     INT            NOT NULL,                      -- 商家ID（关联商家表）
    status        VARCHAR(20)    NOT NULL,                      -- 状态：pending/shipped/out_of_stock/cancelled
    amount        NUMERIC(10, 2) NOT NULL,                      -- 子订单金额（仅含可发货商品）
    created_at    TIMESTAMP DEFAULT NOW()                       -- 子订单创建时间
);

-- 订单商品关联表, 记录子订单中的商品明细
CREATE TABLE order_items
(
    item_id      SERIAL PRIMARY KEY,                       -- 关联项唯一标识
    sub_order_id INT REFERENCES sub_orders (sub_order_id), -- 关联子订单
    product_id   INT            NOT NULL,                  -- 商品ID
    quantity     INT            NOT NULL,                  -- 购买数量
    price        NUMERIC(10, 2) NOT NULL                   -- 商品单价（快照）
);

