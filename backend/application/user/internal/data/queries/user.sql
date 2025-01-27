create database ecommence;

create schema "users";

create table users.users(
    id SERIAL PRIMARY KEY,
    user_id VARCHAR UNIQUE NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    deleted_at timestamptz DEFAULT now() NOT NULL
);

create index user_id_idx ON users.users(user_id);
