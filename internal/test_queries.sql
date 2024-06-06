
-- GetSnippet
SELECT id, title, content, created, datetime(created, expires) as ends 
FROM snippets
WHERE created <= ends AND id=2;

-- GetAllSnippets
SELECT id, title, content, created, datetime(created, expires) as ends 
FROM snippets
WHERE created <= ends
ORDER BY id DESC
LIMIT 10;

-- InsertSnippet
INSERT INTO snippets (title, content, expires)
VALUES (?, ?, ?);

