CREATE SCHEMA IF NOT EXISTS comments;
SET search_path TO comments;

CREATE TABLE comments.comments
(
    id          BIGINT PRIMARY KEY,
    product_id  UUID      NOT NULL,                               -- 商品 ID
    merchant_id UUID      NOT NULL,                               -- 商家 ID
    user_id     UUID      NOT NULL,                               -- 用户 ID
    score       INT       NOT NULL CHECK (score BETWEEN 1 AND 5), -- 评分, 1 - 5
    content     TEXT      NOT NULL,                               -- 评论内容
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE comments.comments IS '评论表';

CREATE INDEX idx_comments_product ON comments.comments (merchant_id, product_id);
CREATE INDEX idx_comments_user ON comments.comments (user_id);