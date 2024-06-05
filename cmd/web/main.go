package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// See Journal, 2024-06-04 19:09 for documentation
type Application struct {
	ErrLog  *log.Logger
	InfoLog *log.Logger
}

func main() {
	// See Journal, 2024-06-04 08:05 for documentation
	port := flag.String(`port`, `:3000`, "setting the port number")
	flag.Parse()

	// see Journal: 2024-06-04 for documentation
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// introduce the app Object in order to grant access to the global
	// application state.
	app := &Application{
		ErrLog:  errLog,
		InfoLog: infoLog,
	}

	// create new servemux (router) where all Routing is initialized.
	mux := app.Routes()

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
