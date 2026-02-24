-- +goose Up
CREATE TABLE cards (
  fragrantica_id INTEGER PRIMARY KEY REFERENCES fragrances(fragrantica_id),
  url TEXT NOT NULL,
  image TEXT NOT NULL,
  has_card BOOLEAN NOT NULL DEFAULT 'f',
  updated TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS cards;