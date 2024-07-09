
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
RETURNING id, title, created;
	
-- name: InsertUser :one
INSERT INTO users (name, email, hashed_password)
VALUES (?, ?, ?)
RETURNING id, name, email, created;

-- name: GetUser :one
SELECT id, name, hashed_password FROM users
WHERE email = ?;

-- name: CheckUser :one
SELECT COUNT(*) from users WHERE id=?;
