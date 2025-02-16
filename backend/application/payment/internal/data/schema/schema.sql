CREATE TABLE table_name
(
    id         SERIAL PRIMARY KEY,
    created_at timestamptz DEFAULT now() NOT NULL,
    deleted_at timestamptz DEFAULT now() NOT NULL
);
