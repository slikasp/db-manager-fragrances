-- name: AddCard :one
INSERT INTO cards (fragrantica_id, url, image, found, downloaded)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: UpdateCard :one
UPDATE cards 
SET image = $2, found = $3, downloaded = $4
WHERE fragrantica_id = $1
RETURNING *;