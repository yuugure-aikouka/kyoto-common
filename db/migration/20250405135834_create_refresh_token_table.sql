-- +goose Up
-- +goose StatementBegin
CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    token TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE refresh_tokens;
-- +goose StatementEnd
