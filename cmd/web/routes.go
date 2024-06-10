package main

import "net/http"

func (app *Application) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	// see Journal: 2024-06-04 for documentation
	fileServer := http.FileServer(http.Dir(`./ui/static/`))

	// Register the fileServer for all URL paths that start with '/static/'.
	// For matching paths, we strip the '/static' prefix before the request
	// reaches the fileServer.
	mux.HandleFunc(`GET /home`, app.handleHome)
	mux.Handle(`/static/`, http.StripPrefix(`/static`, fileServer))

	// Endpoints with handlers as app methods
	mux.HandleFunc(`GET /urlquery`, app.handleUrlQuery)

	mux.HandleFunc(`GET /snippets/`, app.handleSnippetList)
	mux.HandleFunc(`GET /snippets/{id}`, app.handleSingleSnippetView)
	mux.HandleFunc(`POST /snippets/new`, app.handleNewSnippet)

	return mux
}

// vim: ts=4 sw=4 fdm=indent
