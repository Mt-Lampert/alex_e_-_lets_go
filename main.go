package main

import (
	"fmt"
	"log"
	"net/http"
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

func main() {
	// use the http.NewServeMux() constructor to initialize a new servemux (router),
	// then register the home() function as handler for the `/` endpoint.
	mux := http.NewServeMux()
	// This is how it's done in go 1.22+
	mux.HandleFunc(`GET /`, handleHome)

	mux.HandleFunc(`GET /snippets/{id}`, handleSingleSnippetView)
	mux.HandleFunc(`POST /new`, handleNewSnippet)

	// Use the http.ListenAndServe() function as web serving unit. It accepts two parameters:
	//   - the URL (which will be `localhost:3000` here)
	//   - the router we just created.
	// If the webserver returns an error, we handle it using log.Fatal() to log the error and exit.
	// Note that any error returned by http.ListenAndServe() is non-nil!
	log.Println("starting server at port :3000")
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatalf("Uh oh! %s", err)
	}
}

// vim: ts=4 sw=4 fdm=indent
