-- +goose Up
CREATE TABLE perfumers (
  name         TEXT PRIMARY KEY,
  country      TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS perfumers;