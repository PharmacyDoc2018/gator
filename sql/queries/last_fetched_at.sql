-- name: MarkFeedFetched :exec
UPDATE feeds
SET 
    last_fetched_at = $2,
    updated_at = $2
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;