-- +goose Up
CREATE TABLE cards (
  fragrantica_id INTEGER PRIMARY KEY,
  url TEXT NOT NULL,
  image TEXT,
  found BOOLEAN NOT NULL DEFAULT 'f',
  downloaded TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS cards;