-- +goose Up
CREATE TABLE fragrances (
  id BIGSERIAL PRIMARY KEY,

  url          TEXT,
  name         TEXT,
  brand        TEXT,
  country      TEXT,
  gender       TEXT,

  rating_value TEXT,
  rating_count INTEGER,
  year         INTEGER,

  top_notes    TEXT,
  middle_notes TEXT,
  base_notes   TEXT,

  perfumer1    TEXT,
  perfumer2    TEXT,

  accord1      TEXT,
  accord2      TEXT,
  accord3      TEXT,
  accord4      TEXT,
  accord5      TEXT
);

-- +goose Down
DROP TABLE IF EXISTS fragrances;