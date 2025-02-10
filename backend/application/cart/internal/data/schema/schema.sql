CREATE SCHEMA IF NOT EXISTS cart_schema;

CREATE TABLE IF NOT EXISTS cart_schema.carts (
    user_id UUID PRIMARY KEY,
    items JSONB NOT NULL DEFAULT '[]'::jsonb,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);