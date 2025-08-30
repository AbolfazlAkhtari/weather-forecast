-- +goose Up
-- +goose StatementBegin
CREATE TABLE weathers
(
    id          UUID PRIMARY KEY,
    city_name   VARCHAR(255)     NOT NULL,
    country     VARCHAR(255)     NOT NULL,
    temperature DOUBLE PRECISION NOT NULL,
    description VARCHAR(255),
    humidity    INT              NOT NULL,
    wind_speed  DOUBLE PRECISION NOT NULL,
    fetched_at  TIMESTAMP        NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS weathers;
-- +goose StatementEnd
