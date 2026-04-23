-- +goose Up
ALTER TABLE fragrances
ADD COLUMN accord6 TEXT,
ADD COLUMN accord7 TEXT,
ADD COLUMN accord8 TEXT,
ADD COLUMN accord9 TEXT,
ADD COLUMN accord10 TEXT;

-- +goose Down
ALTER TABLE fragrances
DROP COLUMN accord6,
DROP COLUMN accord7,
DROP COLUMN accord8,
DROP COLUMN accord9,
DROP COLUMN accord10;