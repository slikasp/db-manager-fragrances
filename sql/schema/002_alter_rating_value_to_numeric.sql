-- +goose Up
ALTER TABLE fragrances
ALTER COLUMN rating_value
TYPE NUMERIC(3,2)
USING
  NULLIF(REPLACE(rating_value, ',', '.'), '')::NUMERIC(3,2);

-- +goose Down
ALTER TABLE fragrances
ALTER COLUMN rating_value
TYPE TEXT
USING
  CASE
    WHEN rating_value IS NULL THEN NULL
    ELSE TO_CHAR(rating_value, 'FM999999990.00')
  END;