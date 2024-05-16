-- +goose Up
-- +goose StatementBegin
ALTER TABLE url RENAME COLUMN url to fullurl
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE fullurl RENAME COLUMN fullurl to url
-- +goose StatementEnd
