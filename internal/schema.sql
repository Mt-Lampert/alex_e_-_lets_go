-- CREATE TABLE IF NOT EXISTS snippets (
-- 	id INTEGER NOT NULL PRIMARY KEY,
-- 	title VARCHAR(100) NOT NULL,
-- 	content TEXT NOT NULL,
-- 	created DATETIME DEFAULT current_timestamp,
-- 	expires VARCHAR(10) DEFAULT '1 month'
-- );

-- Add an index for the 'created' column
-- CREATE INDEX IF NOT EXISTS snippets_created ON snippets(created);

CREATE TABLE IF NOT EXISTS users (
	id INTEGER NOT NULL PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	email VARCHAR(100) NOT NULL UNIQUE,
	hashed_password CHAR(60) NOT NULL,
	created DATETIME DEFAULT current_timestamp
);

-- vim: ts=4 sw=4 fdm=indent
