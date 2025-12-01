-- +goose Up
-- +goose StatementBegin
ALTER TABLE urls
ADD user_id integer[];
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE urls
DROP COLUMN user_id;
-- +goose StatementEnd
