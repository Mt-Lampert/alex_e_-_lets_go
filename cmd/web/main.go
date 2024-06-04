package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// use the http.NewServeMux() constructor to initialize a new servemux (router),
	// then register the home() function as handler for the `/` endpoint.
	mux := http.NewServeMux()

	// Define a new command line flag with the name 'addr' and a default value
	// of ':3000' and a short help text to tell what this flag is doing.
	port := flag.String(`port`, `:3000`, "setting the port number")

	// Now we have to use the flag.Parse() function to parse the command-line flag.
	// This reads in the command line flag value and assigns it to the 'port' variable.
	// We need to call this **before** we use the 'port' variable; otherwise the value
	// will always be ':3000'.
	// If any errors occur, the application will panic.
	flag.Parse()

	// Create a file server that serves static files out of './ui/static/'. The
	// path here is relative to the project directory root.
	fileServer := http.FileServer(http.Dir(`./ui/static/`))

	// Register the fileServer for all URL paths that start with '/static/'.
	// For matching paths, we strip the '/static' prefix before the request
	// reaches the fileServer.
	http.Handle(`/static/`, http.StripPrefix(`/static`, fileServer))

	// Endpoints
	mux.HandleFunc(`GET /`, handleHome)
	mux.HandleFunc(`GET /urlquery`, handleUrlQuery)

	mux.HandleFunc(`GET /snippets/{id}`, handleSingleSnippetView)
	mux.HandleFunc(`POST /snippets/new`, handleNewSnippet)

	// The value returned from flag.String() is a pointer to the flag value,
	// not the value itself. So we need to dereference the pointer. To make
	// this work properly, Println() must become Printf()
	log.Printf("starting server at port %s", *port)
	err := http.ListenAndServe(*port, mux)
	if err != nil {
		log.Fatalf("Uh oh! %s", err)
	}
}

// vim: ts=4 sw=4 fdm=indent
