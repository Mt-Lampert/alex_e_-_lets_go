package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	// use the http.NewServeMux() constructor to initialize a new servemux (router),
	// then register the home() function as handler for the `/` endpoint.
	mux := http.NewServeMux()

	// See Journal, 2024-06-04 08:05 for documentation
	port := flag.String(`port`, `:3000`, "setting the port number")
	flag.Parse()

	// see Journal: 2024-06-04 for documentation
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// see Journal: 2024-06-04 for documentation
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

	// See 2024-06-04 09:44 Journal for documentation
	srv := &http.Server{
		Addr:     *port,
		ErrorLog: errLog,
		Handler:  mux,
	}

	infoLog.Printf("starting server at port %s", *port)
	// now call the 'ListenAndServe()' method of our own http.Server version
	err := srv.ListenAndServe()
	if err != nil {
		errLog.Fatalf("Uh oh! %s", err)
	}
}

// vim: ts=4 sw=4 fdm=indent
