-- +goose Up
-- +goose StatementBegin
ALTER TABLE urls
ADD users text[];
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE urls
DROP COLUMN users;
-- +goose StatementEnd
