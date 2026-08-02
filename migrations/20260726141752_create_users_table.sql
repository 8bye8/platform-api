-- +goose Up
-- +goose StatementBegin

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
   username VARCHAR(50) NOT NULL UNIQUE,
   email VARCHAR(100) NOT NULL UNIQUE,
   password_hash BYTEA NOT NULL,
   password_salt BYTEA NOT NULL,
   created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX users_credentials_idx
    ON users (username, password_hash, password_salt);

CREATE TRIGGER update_users_updated_at_col
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Dropping the table automatically drops the trigger attached to it
DROP TABLE IF EXISTS users;
-- We drop the function after dropping the table, since the trigger depends on it
DROP FUNCTION IF EXISTS set_updated_at();

-- +goose StatementEnd
