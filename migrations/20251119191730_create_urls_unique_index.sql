-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX unique_urls ON urls (url);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS unique_urls;
-- +goose StatementEnd
