-- +goose Up
-- +goose StatementBegin
CREATE TABLE "user" (
    id uuid PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT uq_email UNIQUE (email)
);

CREATE TABLE token (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL,
    access_token TEXT NOT NULL,
    access_token_expires_at TIMESTAMP NOT NULL,
    refresh_token TEXT NOT NULL,
    refresh_token_expires_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down