-- name: GetFragrance :one
SELECT * 
FROM fragrances 
WHERE fragrantica_id = $1;

-- name: GetFragranceLink :one
SELECT url
FROM fragrances
WHERE fragrantica_id = $1;

-- name: GetFragrancesWithoutDetails :many
SELECT fragrantica_id
FROM fragrances
WHERE name IS NULL;

-- name: AddFragranceLink :exec
INSERT INTO fragrances (fragrantica_id, url, updated)
VALUES (
    $1,
    $2,
    NOW()
);

-- name: UpdateFragrance :exec
UPDATE fragrances 
SET name = $2, brand = $3, country = $4, gender = $5, rating_value = $6, rating_count = $7, year = $8, top_notes = $9, middle_notes = $10, base_notes = $11, perfumer1 = $12, perfumer2 = $13, accord1 = $14, accord2 = $15, accord3 = $16, accord4 = $17, accord5 = $18, updated = NOW()
WHERE fragrantica_id = $1
RETURNING *;
