-- name: GetMissingPerfumers :many
SELECT DISTINCT fragrances.brand
FROM fragrances
WHERE name IS NOT NULL
  AND NOT EXISTS (
      SELECT perfumers.name
      FROM perfumers
      WHERE fragrances.brand = perfumers.name
  );