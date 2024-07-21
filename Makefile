
db_test:
	cp ./snippets_fixtures.db ./snippets.db

db_build:
	cp ./snippets_public.db ./snippets.db

gen_templ:
	templ generate

dev: db_build gen_templ
	go run ./cmd/web/

test: db_test
	go test -v ./cmd/web/



# vim: ts=4 sw=4 fdm=indent
