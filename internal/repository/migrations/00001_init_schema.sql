-- +goose Up
-- +goose StatementBegin
CREATE TABLE restaurants (
    id          UUID PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    logo        TEXT         NOT NULL DEFAULT '',
    active      BOOLEAN      NOT NULL DEFAULT FALSE
);

CREATE TABLE menus (
    id            UUID PRIMARY KEY,
    restaurant_id UUID    NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    description   TEXT    NOT NULL DEFAULT '',
    active        BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE dishes (
    id             SERIAL PRIMARY KEY,
    menu_id        UUID           NOT NULL REFERENCES menus(id) ON DELETE CASCADE,
    name           VARCHAR(255)   NOT NULL,
    compound       TEXT           NOT NULL DEFAULT '',
    price          NUMERIC(10,2)  NOT NULL DEFAULT 0,
    category       VARCHAR(100)   NOT NULL DEFAULT '',
    calories       REAL           NOT NULL DEFAULT 0,
    proteins       REAL           NOT NULL DEFAULT 0,
    fats           REAL           NOT NULL DEFAULT 0,
    carbohydrates  REAL           NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS dishes;
DROP TABLE IF EXISTS menus;
DROP TABLE IF EXISTS restaurants;
-- +goose StatementEnd
