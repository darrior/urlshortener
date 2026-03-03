-- +goose Up
-- +goose StatementBegin
ALTER TABLE urls
ADD deleted boolean;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE urls
DROP COLUMN deleted;
-- +goose StatementEnd
