-- name: AddFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT feeds.name, url, users.name FROM feeds
INNER JOIN users
ON feeds.user_id = users.id;

-- name: GetFeed :one
SELECT * FROM feeds 
WHERE url = $1 LIMIT 1;

-- name: GetFeedsOwned :many
SELECT feeds.name FROM feeds
WHERE feeds.user_id = $1;

-- name: IsOwnerFeed :one
SELECT EXISTS (
    SELECT 1 FROM feeds
    WHERE user_id = $1
    AND feeds.url = $2
);