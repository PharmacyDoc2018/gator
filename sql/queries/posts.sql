-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1, 
    $2, 
    $3, 
    $4, 
    $5, 
    $6, 
    $7, 
    $8
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.title, posts.url, posts.description, posts.published_at 
FROM posts
LEFT JOIN feeds ON posts.feed_id = feeds.id
LEFT JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
WHERE feeds.user_id = $1
OR feed_follows.user_id = $1
ORDER BY published_at DESC 
LIMIT $2;
