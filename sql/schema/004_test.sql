-- +goose Up
ALTER TABLE fragrances
ADD COLUMN test TEXT;

-- +goose Down
ALTER TABLE fragrances
DROP COLUMN test;