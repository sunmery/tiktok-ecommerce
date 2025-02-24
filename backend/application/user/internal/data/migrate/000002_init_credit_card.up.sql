CREATE SCHEMA IF NOT EXISTS users;
SET search_path TO users;

CREATE TABLE users.credit_cards
(
    id               SERIAL PRIMARY KEY,
    user_id          UUID        NOT NULL,
    number           VARCHAR(20) NOT NULL,
    cvv              VARCHAR(4)  NOT NULL,
    expiration_year  CHAR(4)     NOT NULL,
    expiration_month CHAR(2)     NOT NULL
);

CREATE INDEX idx_credit_cards_id ON users.credit_cards (id);
