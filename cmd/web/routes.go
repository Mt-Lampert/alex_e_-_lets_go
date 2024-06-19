package main

import "net/http"

// Single source of truth for routing in this app
func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	// see Journal: 2024-06-04 for documentation
	fileServer := http.FileServer(http.Dir(`./ui/static/`))

	// Register the fileServer for all URL paths that start with '/static/'.
	// For matching paths, we strip the '/static' prefix before the request
	// reaches the fileServer.
	mux.Handle(`/static/`, http.StripPrefix(`/static`, fileServer))

	mux.HandleFunc(`GET /home`, app.handleHome)
	// Endpoints with handlers as app methods
	mux.HandleFunc(`GET /urlquery`, app.handleUrlQuery)

	mux.HandleFunc(`GET /snippets/`, app.handleSnippetList)
	mux.HandleFunc(`GET /snippets/{id}`, app.handleSingleSnippetView)
	mux.HandleFunc(`POST /snippets/new`, app.handleNewSnippet)

	mwChain := createMdwChain(
		app.recoverPanic,
		app.logRequest,
		secureHeaders,
	)

	return mwChain(mux)
}

// vim: ts=4 sw=4 fdm=indent
