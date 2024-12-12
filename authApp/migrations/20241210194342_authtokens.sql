-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE tokens (
    user_id VARCHAR(255) NOT NULL,
    token VARCHAR(500) NOT NULL,
    PRIMARY KEY (token)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS tokens;
-- +goose StatementEnd