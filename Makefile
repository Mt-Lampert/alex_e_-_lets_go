
db_test:
	cp ./snippets_fixtures.db ./snippets.db

db_build:
	cp ./snippets_public.db ./snippets.db


dev: db_build
	go run ./cmd/web/

test: db_test
	go test -v ./cmd/web/

