-- +goose Up
ALTER TABLE fragrances
ADD COLUMN updated TIMESTAMP NOT NULL DEFAULT NOW();

-- +goose Down
ALTER TABLE fragrances
DROP COLUMN updated;