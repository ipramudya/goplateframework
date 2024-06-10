-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS addresses;
DROP INDEX IF EXISTS addresses_idx;

CREATE TABLE addresses (
    id          uuid PRIMARY KEY            NOT NULL   DEFAULT gen_random_uuid(),
    street      varchar(255)                NOT NULL,
    city        varchar(50)                 NOT NULL,
    province    varchar(50)                 NOT NULL,
    postal_code varchar(10)                 NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE    NOT NULL   DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE    NOT NULL   DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX addresses_idx ON addresses(street, city, province, postal_code);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS addresses;
DROP INDEX IF EXISTS addresses_idx;
-- +goose StatementEnd
