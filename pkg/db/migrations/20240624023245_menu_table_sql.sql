-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS menus;
DROP INDEX IF EXISTS menus_name_idx;

CREATE TABLE IF NOT EXISTS
    menus (
        id              uuid PRIMARY KEY            NOT NULL    DEFAULT gen_random_uuid(),
        name            varchar(50)                 NOT NULL,
        description     varchar(255)                NOT NULL,
        price           numeric(10,2)               NOT NULL,
        is_available    boolean                     NOT NULL    DEFAULT false,
        image_url       varchar(255)                NULL        DEFAULT '',
        created_at      TIMESTAMP WITH TIME ZONE    NOT NULL    DEFAULT CURRENT_TIMESTAMP,
        updated_at      TIMESTAMP WITH TIME ZONE    NOT NULL    DEFAULT CURRENT_TIMESTAMP,
        outlet_id       uuid                        NOT NULL,
        
        FOREIGN KEY (outlet_id) REFERENCES outlets(id) ON DELETE CASCADE
    );
CREATE INDEX IF NOT EXISTS menus_name_idx ON menus (name, created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS menus;
DROP INDEX IF EXISTS menus_name_idx;
-- +goose StatementEnd
