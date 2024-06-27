
-- name: GetSnippet :one
SELECT id, title, content, created, datetime(created, expires) as ends 
FROM snippets
WHERE created <= ends AND id=?;

-- name: GetAllSnippets :many
SELECT id, title, content, created, expires
FROM snippets
ORDER BY id DESC
LIMIT 10;

-- name: InsertSnippet :one
INSERT INTO snippets (title, content, expires)
VALUES (?, ?, ?)
RETURNING id, title created;
	
