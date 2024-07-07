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

	// Use the general middleware
	mux.Use(app.recoverPanic)
	mux.Use(app.logRequest)
	mux.Use(secureHeaders)

	// see Journal: 2024-06-04 for documentation

	// Register the fileServer for all URL paths that start with '/static/'.
	// For matching paths, we strip the '/static' prefix before the request
	// reaches the fileServer.
	fileServer := http.FileServer(http.Dir(`./ui/static`))
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// define a new subgroup with its own sub-router 'r'
	// and its own middleware
	mux.Group(func(r chi.Router) {
		r.Use(app.sessionManager.LoadAndSave)
		r.Use(noSurf)

		// Endpoints with handlers as app methods
		r.Get(`/`, app.handleHome)
		r.Get(`/snippets/{id}`, app.handleSingleSnippetView)
		r.Get(`/snippets`, app.handleSnippetList)
		r.Get(`/urlquery`, app.handleUrlQuery)
		r.Get(`/user/login_form`, app.handleLoginForm)
		r.Post(`/user/login`, app.handleLogin)
		r.Get(`/user/signup_form`, app.handleSignupForm)
		r.Post(`/user/signup`, app.handleSignup)

		r.Group(func(rAuth chi.Router) {
			rAuth.Use(app.requireAuthentication)

			rAuth.Get(`/new/snippet`, app.handleNewSnippetForm)
			rAuth.Post(`/create/snippet`, app.handleNewSnippet)
			rAuth.Post(`/user/logout`, app.handleLogout)
		})
	})

	return mux
}

// vim: ts=4 sw=4 fdm=indent
