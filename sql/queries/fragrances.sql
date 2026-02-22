-- name: GetFragrance :one
SELECT * 
FROM fragrances 
WHERE fragrantica_id = $1;