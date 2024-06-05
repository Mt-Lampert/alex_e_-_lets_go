package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *Application) handleHome(w http.ResponseWriter, r *http.Request) {
	//exclude anything but root as endpoint
	if r.URL.Path != `/` {
		http.NotFound(w, r)
		return
	}

	templates := []string{
		"./ui/html/base.go.html",
		"./ui/html/pages/home.go.html",
		"./ui/html/partials/nav.go.html",
	}

	// See Journal, 2024-06-03 07:43 for documentation
	ts, err := template.ParseFiles(templates...)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// See Journal, 2024-06-03 07:43 for documentation
	err = ts.ExecuteTemplate(w, `base`, nil)
	if err != nil {
		app.ServerError(w, err)
		return
	}
}

// Add a handler function for creating a snippet.
func (app Application) handleNewSnippet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `Creating a new snippet ...`)
}

// Add a handler function for viewing a specific snippet
func (app Application) handleSingleSnippetView(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue(`id`)
	fmt.Fprintf(w, `Display snippet with ID '%s'`, id)
}

func (app Application) handleUrlQuery(w http.ResponseWriter, r *http.Request) {
	// get the id
	rawId := r.URL.Query().Get(`id`)

	// validate the id:
	//    - it must be numerical
	//    - it must be greater than 0
	id, err := strconv.Atoi(rawId)
	if err != nil || id <= 0 {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, `Seems you are looking for something with id '%s'`, rawId)
}

// vim: ts=4 sw=4 fdm=indent
