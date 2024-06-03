-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS accounts;
DROP TYPE IF EXISTS role;
DROP INDEX IF EXISTS email_idx;

CREATE TYPE role AS ENUM('user', 'admin', 'superadmin');

CREATE TABLE IF NOT EXISTS
   accounts (
      id          uuid PRIMARY KEY            NOT NULL   DEFAULT gen_random_uuid(),
      firstname   varchar(50)                 NOT NULL,
      lastname    varchar(50)                 NOT NULL,
      email       varchar(50) UNIQUE          NOT NULL,
      password    varchar(255)                NOT NULL,
      phone       varchar(50)                 NOT NULL,
      role        role                        NOT NULL   DEFAULT 'user',
      created_at  TIMESTAMP WITH TIME ZONE    NOT NULL   DEFAULT CURRENT_TIMESTAMP,
      updated_at  TIMESTAMP WITH TIME ZONE    NOT NULL   DEFAULT CURRENT_TIMESTAMP
   );

CREATE UNIQUE INDEX IF NOT EXISTS email_idx ON accounts (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS accounts;

DROP INDEX IF EXISTS email_idx;

DROP TYPE IF EXISTS role;
-- +goose StatementEnd
