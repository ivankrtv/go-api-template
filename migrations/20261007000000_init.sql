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
-- +goose StatementEnd

-- +goose Down