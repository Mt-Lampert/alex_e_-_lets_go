package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// writes an error message and stack trace to the error log.
// Then it sends a generic 500 Internal Server Error message to the frontend.
func (app *Application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf(`%s\n%s`, err.Error(), debug.Stack())
	app.ErrLog.Println(trace)

	http.Error(w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}

// sends a specific status code and corresponding description to the user.
// This is for messages like "400 Bad Request".
func (app *Application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// For consistency, we build a "Not Found" helper. This is simply a convenient
// wrapper around 'clientError()' which sends a "404 Not Found" error to the
// frontend.
func (app *Application) NotFound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}

// vim: ts=4 sw=4 fdm=indent
