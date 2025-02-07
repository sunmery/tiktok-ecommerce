CREATE TABLE IF NOT EXISTS products.products
(
    id           SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description  TEXT        NOT NULL,
    picture      TEXT        NOT NULL,
    price        REAL        NOT NULL,
    categories   TEXT[]
);
