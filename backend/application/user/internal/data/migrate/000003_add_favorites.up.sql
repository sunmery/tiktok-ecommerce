CREATE SCHEMA IF NOT EXISTS users;
SET search_path TO users;

CREATE TABLE users.favorites
(
    user_id    UUID NOT NULL,
    product_id UUID NOT NULL,
    merchant_id UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    PRIMARY KEY (user_id, product_id)
);
COMMENT ON TABLE users.favorites IS '消费者收藏表';

CREATE INDEX idx_users_favorites ON users.favorites(user_id, merchant_id);
