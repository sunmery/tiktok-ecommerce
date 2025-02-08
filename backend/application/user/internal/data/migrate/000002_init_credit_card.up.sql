CREATE SCHEMA credit_cards;

CREATE TABLE credit_cards.credit_cards
(
    id               SERIAL PRIMARY KEY,
    owner            VARCHAR(100) NOT NULL DEFAULT 'tiktok',
    name             VARCHAR(100) NOT NULL,
    number           VARCHAR(20)  NOT NULL,
    cvv              VARCHAR(4)   NOT NULL,
    expiration_year  CHAR(4)      NOT NULL,
    expiration_month CHAR(2)      NOT NULL
);

CREATE INDEX idx_credit_cards_id ON credit_cards.credit_cards (id);
