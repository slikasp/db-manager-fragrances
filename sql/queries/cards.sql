-- name: AddCard :one
INSERT INTO cards (fragrantica_id, url, image, has_card, updated)
VALUES (
    $1,
    $2,
    $3,
    $4,
    NOW()
)
RETURNING *;

-- name: UpdateCard :one
UPDATE cards 
SET image = $2, has_card = $3, updated = NOW()
WHERE fragrantica_id = $1
RETURNING *;

-- name: RefreshCard :exec
UPDATE cards 
SET downloaded = NOW()
WHERE fragrantica_id = $1;