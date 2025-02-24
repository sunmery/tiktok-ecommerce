CREATE SCHEMA IF NOT EXISTS users;
SET search_path TO users;

CREATE TABLE users.addresses
(
    id             SERIAL PRIMARY KEY,
    user_id        UUID         NOT NULL,
    street_address TEXT         NOT NULL, -- 街道地址
    city           VARCHAR(255) NOT NULL, -- 城市
    state          VARCHAR(20)  NOT NULL, -- 状态
    country        VARCHAR(100) NOT NULL, -- 国家
    zip_code       VARCHAR(20)  NOT NULL  -- 邮政编码
);

CREATE TABLE users.credit_cards
(
    id               SERIAL PRIMARY KEY,
    user_id          UUID        NOT NULL,
    number           VARCHAR(20) NOT NULL,
    cvv              VARCHAR(4)  NOT NULL,
    expiration_year  CHAR(4)     NOT NULL,
    expiration_month CHAR(2)     NOT NULL
);

