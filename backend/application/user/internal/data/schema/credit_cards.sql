CREATE SCHEMA IF NOT EXISTS users;
SET search_path TO users;

CREATE TABLE users.credit_cards
(
    id           SERIAL PRIMARY KEY,                 -- 卡 ID
    user_id      uuid         NOT NULL,              --用户ID
    number       VARCHAR(20)  NOT NULL,              -- 卡号
    currency     VARCHAR(3)   NOT NULL,              -- 货币类型
    cvv          VARCHAR(4)   NOT NULL,              -- 安全码
    exp_year     VARCHAR(4)   NOT NULL,              -- 过期年份
    exp_month    VARCHAR(2)   NOT NULL,              -- 过期月份
    owner        VARCHAR(100) NOT NULL,              -- 持卡人姓名
    name         VARCHAR(100),                       -- 卡名
    type         VARCHAR(20)  NOT NULL,              -- 卡类型（如借记卡、信用卡）
    brand        VARCHAR(20)  NOT NULL,              -- 卡品牌（如 Visa、MasterCard）
    country      VARCHAR(50)  NOT NULL,              -- 卡所属国家
    created_time timestamptz    NOT NULL DEFAULT NOW() -- 创建时间
);
