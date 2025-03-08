-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN email TEXT UNIQUE NOT NULL;

CREATE TABLE magic_links (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    token TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS magic_links;

ALTER TABLE users DROP COLUMN IF EXISTS email;
-- +goose StatementEnd
