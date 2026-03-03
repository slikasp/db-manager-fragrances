-- name: GetMissingFragranceIDs :many
SELECT cards.fragrantica_id
FROM cards
WHERE has_card = 't'
  AND NOT EXISTS (
      SELECT fragrances.fragrantica_id
      FROM fragrances
      WHERE cards.fragrantica_id = fragrances.fragrantica_id
  )
ORDER BY cards.fragrantica_id ASC;