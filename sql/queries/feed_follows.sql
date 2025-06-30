-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT (
    inserted_feed_follow.*,
    feeds.name,
    users.name
)
FROM inserted_feed_follow
INNER JOIN users 
ON inserted_feed_follow.user_id = users.id
INNER JOIN feeds
ON inserted_feed_follow.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT feeds.name, users.name 
FROM feed_follows
INNER JOIN feeds
ON feed_follows.feed_id = feeds.id
INNER JOIN users
ON feeds.user_id = users.id
WHERE feed_follows.user_id = $1;

-- name: IsFollowingFeed :one
SELECT EXISTS (
    SELECT 1 FROM feed_follows
    INNER JOIN feeds 
    ON feed_follows.feed_id = feeds.id
    WHERE feed_follows.user_id = $1
    AND feeds.url = $2
);

-- name: DeleteFollow :exec
DELETE FROM feed_follows
USING feeds
WHERE feed_follows.feed_id = feeds.id
AND feed_follows.user_id = $1
AND feeds.url = $2;
