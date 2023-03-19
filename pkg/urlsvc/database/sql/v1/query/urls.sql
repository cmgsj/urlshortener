-- name: GetUrl :one
SELECT * FROM urls WHERE url_id = ? LIMIT 1;

-- name: ListUrls :many
SELECT * FROM urls;

-- name: CreateUrl :one
INSERT INTO urls (url_id, redirect_url) VALUES (?, ?) RETURNING *;

-- name: UpdateUrl :exec
UPDATE urls SET redirect_url = ? WHERE url_id = ?;

-- name: DeleteUrl :exec
DELETE FROM urls WHERE url_id = ?;