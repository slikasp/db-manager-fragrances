-- name: GetFragrance :one
SELECT * 
FROM fragrances 
WHERE fragrantica_id = $1;

-- name: GetFragranceLink :one
SELECT url
FROM fragrances
WHERE fragrantica_id = $1;

-- name: AddFragranceLink :exec
INSERT INTO fragrances (fragrantica_id, url, updated)
VALUES (
    $1,
    $2,
    NOW()
);