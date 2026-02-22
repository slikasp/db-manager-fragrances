-- +goose Up
ALTER TABLE fragrances
ADD COLUMN fragrantica_id integer;

UPDATE fragrances
SET fragrantica_id = substring(url FROM '-([0-9]+)\.html$')::integer;

ALTER TABLE fragrances
ALTER COLUMN fragrantica_id SET NOT NULL;

ALTER TABLE fragrances
ADD CONSTRAINT fragrances_fragrantica_id_key UNIQUE (fragrantica_id);


-- +goose Down
ALTER TABLE fragrances
DROP CONSTRAINT fragrances_fragrantica_id_key;

ALTER TABLE fragrances
DROP COLUMN fragrantica_id;