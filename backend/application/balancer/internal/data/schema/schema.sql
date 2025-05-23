CREATE SCHEMA IF NOT EXISTS balances;

CREATE TABLE balances.user_balances
(
    user_id    UUID UNIQUE    NOT NULL,
    currency   CHAR(3),
    available  DECIMAL(12, 2) NOT NULL DEFAULT 0.00 CHECK (available >= 0), -- 可用余额
    frozen     DECIMAL(12, 2) NOT NULL DEFAULT 0.00 CHECK (frozen >= 0),    -- 冻结余额
    version    INT            NOT NULL DEFAULT 0,-- 乐观锁
    created_at timestamptz             DEFAULT NOW(),
    updated_at timestamptz             DEFAULT NOW(),
    PRIMARY KEY (user_id, currency)
);
COMMENT ON TABLE balances.user_balances IS '消费者余额表';

-- 用户支付方式表
CREATE TABLE balances.user_payment_methods
(
    id              BIGINT PRIMARY KEY,
    user_id         UUID        NOT NULL,
    type            VARCHAR(20) NOT NULL
        CHECK ( type IN ('ALIPAY', 'BANK_CARD', 'BALANCE')),
    is_default      BOOLEAN     NOT NULL DEFAULT FALSE,
    account_details JSONB       NOT NULL DEFAULT '{}', -- 示例: {"account": "123@alipay.com", "real_name": "张三"}
    created_at      timestamptz NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE balances.user_payment_methods IS '用户支付方式表';

-- 冻结记录表
CREATE TABLE balances.balance_freezes
(
    id         BIGINT PRIMARY KEY,
    user_id    UUID           NOT NULL,
    order_id   BIGINT         NOT NULL, -- 关联订单号
    currency   CHAR(3)        NOT NULL, -- 币种 ISO 4217 代码
    amount     DECIMAL(12, 2) NOT NULL CHECK (amount > 0),
    status     VARCHAR(20)    NOT NULL  -- 冻结状态: 冻结|确认|取消
        CHECK (status IN ('FROZEN', 'CONFIRMED', 'CANCELED')),
    created_at timestamptz    NOT NULL DEFAULT NOW(),
    updated_at timestamptz    NOT NULL DEFAULT NOW(),
    expires_at timestamptz    NOT NULL  -- 冻结过期时间
);
COMMENT ON TABLE balances.balance_freezes IS '冻结记录表';

-- 商家支付方式表 (merchant_payment_methods)
-- 支持商家绑定多个支付账号（支付宝、微信、银行账户）。
-- CREATE TYPE merchant_payment_type AS ENUM ('ALIPAY', 'WECHAT', 'BANK_ACCOUNT','BALANCE');
CREATE TABLE balances.merchant_payment_methods
(
    id              BIGINT PRIMARY KEY,
    merchant_id     UUID NOT NULL,
    type            VARCHAR(15) NOT NULL  -- 支付方式
        CHECK ( type IN ('ALIPAY', 'BANK_CARD', 'BALANCE')),
    is_default      BOOLEAN     NOT NULL DEFAULT FALSE,
    account_details JSONB       NOT NULL, -- 示例: {"account": "merchant@alipay.com", "bank_name": "中国银行"}
    created_at      timestamptz NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE balances.merchant_payment_methods IS '商家支付方式表';

-- 商家余额表
-- 资金要从用户冻结余额转移到商家
CREATE TABLE balances.merchant_balances
(
    merchant_id UUID           NOT NULL,
    currency    CHAR(3)        NOT NULL,
    available   DECIMAL(12, 2) NOT NULL DEFAULT 0 CHECK (available >= 0), -- 可用余额
    version     INT            NOT NULL DEFAULT 0,                        -- 乐观锁
    PRIMARY KEY (merchant_id, currency),
    created_at  timestamptz    NOT NULL DEFAULT NOW(),
    updated_at  timestamptz    NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE balances.merchant_balances IS '商家余额表';

-- 交易流水表 (transactions)
-- 记录所有资金变动，保留支付方式快照
-- CREATE TYPE transaction_type AS ENUM ('RECHARGE', 'PAYMENT', 'REFUND', 'WITHDRAW');
-- CREATE TYPE transaction_status AS ENUM ('PENDING', 'SUCCESS', 'FAILED');
CREATE TABLE balances.transactions
(
    id                  BIGINT PRIMARY KEY,
    type                VARCHAR(15)    NOT NULL                               -- 交易类型: 充值|支付|退款|提现
        CHECK ( type IN ('RECHARGE', 'PAYMENT', 'REFUND', 'WITHDRAW')),
    amount              DECIMAL(12, 2) NOT NULL CHECK (amount > 0),
    currency            CHAR(3)        NOT NULL,                              -- 交易币种 ISO 4217 代码

    -- 交易双方
    from_user_id        UUID           NOT NULL,
    to_merchant_id      UUID           NOT NULL,

    -- 支付方式快照
    payment_method_type VARCHAR(20)    NOT NULL,                              -- 支付方式: ALIPAY|WECHAT|BALANCE|BANK_CARD
    payment_account     VARCHAR(255)   NOT NULL,                              -- 支付账号, 对应的第三方支付方式的账号
    payment_extra       JSONB                   DEFAULT '{}'::JSONB NOT NULL, -- 交易号等额外信息

    status              VARCHAR(15)    NOT NULL DEFAULT 'PENDING'             -- 支付状态: 等待支付|已支付|取消支付|支付异常
        CHECK ( status IN ('PENDING', 'PAID', 'CANCELLED', 'FAILED') ),
    freeze_id           BIGINT         NOT NULL,                              -- 关联冻结记录
    idempotency_key     VARCHAR(255)   NOT NULL,                              -- 幂等键
    consumer_version    BIGINT         NOT NULL,                              -- 用户乐观锁版本
    merchant_version    BIGINT         NOT NULL,                              -- 商家乐观锁版本
    created_at          timestamptz    NOT NULL DEFAULT NOW(),
    updated_at          timestamptz    NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE balances.merchant_balances IS '交易流水表';
