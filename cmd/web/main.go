package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/MtLampert/alex_e_-_lets_go/internal/db"
	"github.com/go-playground/form/v4"
)

// See Journal, 2024-06-04 19:09 for documentation
// Adding a template cache struct
type Application struct {
	ErrLog        *log.Logger
	InfoLog       *log.Logger
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

func main() {
	// initializing the database module
	db.Setup()

	// See Journal, 2024-06-04 08:05 for documentation
	port := flag.String(`port`, `:3000`, "setting the port number")
	flag.Parse()

	// see Journal: 2024-06-04 for documentation
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// build instance of a template cache
	templateCache, err := buildTemplateCache()
	if err != nil {
		errLog.Fatal(err)
	}

	// Initialize a new form decoder
	formDecoder := form.NewDecoder()

	// introduce the app Object in order to grant access to the global
	// application state.
	app := &Application{
		ErrLog:        errLog,
		InfoLog:       infoLog,
		templateCache: templateCache,
		formDecoder:   formDecoder,
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
	err = srv.ListenAndServe()
	if err != nil {
		errLog.Fatalf("Uh oh! %s", err)
	}
}

// vim: ts=4 sw=4 fdm=indent
