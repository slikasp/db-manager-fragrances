-- name: AddPerfumer :one
INSERT INTO perfumers (name, country)
VALUES (
    $1,
    $2
)
RETURNING *;

-- name: GetPerfumerCountry :one
SELECT country
FROM perfumers
WHERE name = $1
LIMIT 1;