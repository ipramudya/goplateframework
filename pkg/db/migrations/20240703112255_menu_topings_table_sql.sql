-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS menu_topings;
DROP INDEX IF EXISTS menu_topings_name_idx;

CREATE TABLE IF NOT EXISTS
    menu_topings (
        id              uuid PRIMARY KEY            NOT NULL    DEFAULT gen_random_uuid(),
        name            varchar(50)                 NOT NULL,
        price           numeric(10,2)               NOT NULL,
        is_available    boolean                     NOT NULL    DEFAULT false,
        image_url       varchar(255)                NULL        DEFAULT '',
        stock           integer                     NOT NULL    DEFAULT 0,
        created_at      TIMESTAMP WITH TIME ZONE    NOT NULL    DEFAULT CURRENT_TIMESTAMP,
        updated_at      TIMESTAMP WITH TIME ZONE    NOT NULL    DEFAULT CURRENT_TIMESTAMP,
        menu_id         uuid                        NOT NULL,
        
        FOREIGN KEY (menu_id) REFERENCES menus(id) ON DELETE CASCADE
    );
CREATE INDEX IF NOT EXISTS menu_topings_name_idx ON menu_topings (name, created_at, updated_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS menu_topings;
DROP INDEX IF EXISTS menu_topings_name_idx;
-- +goose StatementEnd
