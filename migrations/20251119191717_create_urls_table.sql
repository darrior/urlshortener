-- +goose Up
-- +goose StatementBegin
CREATE TABLE urls (id text NOT NULL, url text, PRIMARY KEY(id));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS urls;
-- +goose StatementEnd
