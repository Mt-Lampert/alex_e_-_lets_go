package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Single source of truth for routing in this app
func (app *Application) Routes() *chi.Mux {
	mux := chi.NewRouter()

	// assign a custom http.HandlerFunc as the default handler for 404 Not
	// Found cases.
	// See :GoDoc for documentation
	mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.NotFound(w)
	})

	// Use the middleware
	mux.Use(app.recoverPanic)
	mux.Use(app.logRequest)
	mux.Use(secureHeaders)

	// see Journal: 2024-06-04 for documentation
	fileServer := http.FileServer(http.Dir(`./ui/static`))

	// Register the fileServer for all URL paths that start with '/static/'.
	// For matching paths, we strip the '/static' prefix before the request
	// reaches the fileServer.
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	mux.Get(`/`, app.handleHome)

	// Endpoints with handlers as app methods
	mux.Get(`/urlquery`, app.handleUrlQuery)

	mux.Get(`/snippets`, app.handleSnippetList)
	mux.Get(`/snippets/{id}`, app.handleSingleSnippetView)
	mux.Post(`/snippets/new`, app.handleNewSnippet)

	return mux
}

// vim: ts=4 sw=4 fdm=indent
