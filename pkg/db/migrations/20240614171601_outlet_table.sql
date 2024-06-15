-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS outlets;
DROP INDEX IF EXISTS outlets_name_idx;

DROP FUNCTION IF EXISTS delete_address_on_outlet_delete();
DROP TRIGGER IF EXISTS trigger_delete_address_on_outlet_delete ON outlets;

CREATE TABLE IF NOT EXISTS 
    outlets (
        id              uuid PRIMARY KEY                NOT NULL    DEFAULT gen_random_uuid(),
        name            varchar(50)                     NOT NULL,
        phone           varchar(20)                     NOT NULL,
        opening_time    TIMESTAMP WITH TIME ZONE        NOT NULL,
        closing_time    TIMESTAMP WITH TIME ZONE        NOT NULL,
        created_at      TIMESTAMP WITH TIME ZONE        NOT NULL    DEFAULT CURRENT_TIMESTAMP,
        updated_at      TIMESTAMP WITH TIME ZONE        NOT NULL    DEFAULT CURRENT_TIMESTAMP,
        address_id      uuid                            NOT NULL,
        
        FOREIGN KEY (address_id) REFERENCES addresses(id)
    );
CREATE UNIQUE INDEX IF NOT EXISTS outlets_name_idx ON outlets (name);

CREATE FUNCTION delete_address_on_outlet_delete() RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM addresses WHERE id = OLD.address_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_delete_address_on_outlet_delete
AFTER DELETE ON outlets
FOR EACH ROW
EXECUTE FUNCTION delete_address_on_outlet_delete();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS outlets;
DROP INDEX IF EXISTS outlets_name_idx;
DROP FUNCTION IF EXISTS delete_address_on_outlet_delete();
DROP TRIGGER IF EXISTS trigger_delete_address_on_outlet_delete ON outlets;
-- +goose StatementEnd
