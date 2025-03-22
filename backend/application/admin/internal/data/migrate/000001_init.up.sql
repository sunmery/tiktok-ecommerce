CREATE TABLE table_name
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR                   NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    deleted_at timestamptz DEFAULT now() NOT NULL
);

CREATE INDEX idx_col_name ON table_name (name);
