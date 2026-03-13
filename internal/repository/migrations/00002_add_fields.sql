-- +goose Up
-- +goose StatementBegin
ALTER TABLE restaurants
    ADD COLUMN phone      VARCHAR(50)              DEFAULT '',
    ADD COLUMN address    TEXT                     DEFAULT '',
    ADD COLUMN created_at TIMESTAMPTZ NOT NULL     DEFAULT NOW();

ALTER TABLE menus
    ADD COLUMN name       VARCHAR(255) NOT NULL    DEFAULT '',
    ADD COLUMN created_at TIMESTAMPTZ NOT NULL     DEFAULT NOW();

ALTER TABLE dishes
    ADD COLUMN weight_grams REAL        NOT NULL   DEFAULT 0,
    ADD COLUMN created_at   TIMESTAMPTZ NOT NULL   DEFAULT NOW();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE restaurants
    DROP COLUMN IF EXISTS phone,
    DROP COLUMN IF EXISTS address,
    DROP COLUMN IF EXISTS created_at;

ALTER TABLE menus
    DROP COLUMN IF EXISTS name,
    DROP COLUMN IF EXISTS created_at;

ALTER TABLE dishes
    DROP COLUMN IF EXISTS weight_grams,
    DROP COLUMN IF EXISTS created_at;
-- +goose StatementEnd
