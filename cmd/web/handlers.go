package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
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

	// template.ParseFiles() reads the templates into a template set.
	// If there is an error, we log the detailed error message on the terminal
	// and use the http.Error() function to send a generic 500 server error.
	// the paths are a variadic parameter here.
	ts, err := template.ParseFiles(templates...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, `Internal Server Error has occured!`, http.StatusInternalServerError)
	}

	// Since we made it here, we use the ExecuteTemplate() method on the
	// template set to write the template content as the response body.
	// The second parameter of ExecuteTemplate() is the 'base' template or the
	// master template that this page is built on.
	// The last parameter of Execute() represents any dynamic data we want to
	// pass in; at the moment it will be nil.
	err = ts.ExecuteTemplate(w, `base`, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, `Template Error. WTF?!!`, 500)
	}
}

// Add a handler function for creating a snippet.
func handleNewSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`Creating a new snippet ...`))
}

// Add a handler function for viewing a specific snippet
func handleSingleSnippetView(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue(`id`)
	w.Write([]byte(fmt.Sprintf("Display snippet with ID '%s'", id)))
}

func handleUrlQuery(w http.ResponseWriter, r *http.Request) {
	// get the id
	rawId := r.URL.Query().Get(`id`)

	// validate the id:
	//    - it must be numerical
	//    - it must be greater than 0
	id, err := strconv.Atoi(rawId)
	if err != nil || id <= 0 {
		http.Error(w, `invalid ID!`, http.StatusBadRequest)
		return
	}

	w.Write([]byte(fmt.Sprintf(`Seems you are looking for something with id '%s'`, rawId)))
}

// vim: ts=4 sw=4 fdm=indent
