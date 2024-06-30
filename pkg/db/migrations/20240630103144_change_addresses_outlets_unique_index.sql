-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS addresses_idx;
DROP INDEX IF EXISTS outlets_name_idx;

CREATE INDEX IF NOT EXISTS addresses_idx ON addresses(street, city, province, postal_code);
CREATE INDEX IF NOT EXISTS outlets_name_idx ON outlets (name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS addresses_idx ON addresses(street, city, province, postal_code);
CREATE UNIQUE INDEX IF NOT EXISTS outlets_name_idx ON outlets (name);
-- +goose StatementEnd
