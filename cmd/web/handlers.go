package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// Define a simple home Handler function which writes a byte slice
// containing "Hello from Snippetbox" as a response body
func handleHome(w http.ResponseWriter, r *http.Request) {
	// Write() accepts only []byte as ‘most neutral’ message type
	w.Write([]byte(`Hello from Snippetbox!`))
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
