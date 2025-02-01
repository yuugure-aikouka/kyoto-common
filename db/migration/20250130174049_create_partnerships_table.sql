-- +goose Up
-- +goose StatementBegin
CREATE TYPE partnership_status AS ENUM ('pending', 'accepted', 'rejected');

CREATE TABLE partnerships (
    user_id_1 INT NOT NULL,
    user_id_2 INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id_1, user_id_2),
    FOREIGN KEY (user_id_1) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id_2) REFERENCES users(id) ON DELETE CASCADE,
    status partnership_status NOT NULL DEFAULT 'pending'
);

CREATE INDEX idx_user_id_1 ON partnerships (user_id_1);
CREATE INDEX idx_user_id_2 ON partnerships (user_id_2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE partnerships;
DROP TYPE partnership_status;
-- +goose StatementEnd
