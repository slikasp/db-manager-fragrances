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

-- name: GetCard :one
SELECT *
FROM cards
WHERE fragrantica_id = $1;

-- name: GetLastCardID :one
SELECT fragrantica_id
FROM cards
WHERE has_card = 't'
ORDER BY fragrantica_id DESC
LIMIT 1;

-- name: GetMissingCardIDs :many
SELECT fragrantica_id
FROM cards
WHERE has_card = 'f';

-- name: GetExistingCardIDs :many
SELECT fragrantica_id
FROM cards
WHERE has_card = 't';

-- name: UpdateCard :one
UPDATE cards 
SET image = $2, has_card = $3, updated = NOW()
WHERE fragrantica_id = $1
RETURNING *;

-- name: RefreshCard :exec
UPDATE cards 
SET downloaded = NOW()
WHERE fragrantica_id = $1;