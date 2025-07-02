-- name: MarkFeedFetched :exec
UPDATE feeds
SET 
    updated_at = $2,
    last_fetched_at = $2
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;