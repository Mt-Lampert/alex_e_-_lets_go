package main

import (
	"fmt"
	"net/http"
)

// new type representing a middleware function
type Middleware func(http.Handler) http.Handler

func createMdwChain(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		// building a 'triangle' of nested middleware functions.
		// If you can't get your head around this, follow it step by step in
		// the debugger. No, seriously! You really should understand what's
		// happening here!
		//
		// we are moving bottom-up in the 'xs' list that has been passed to us!
		for i := len(xs) - 1; i >= 0; i-- {
			// xs[i] is the current Middleware function in the list
			x := xs[i]
			// x(next) is the return value of the current Middleware function.
			// This is what creates the nested function calls!
			next = x(next)
		}
		return next
	}
}

// add security for the browser
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

// log all requests
func (app *Application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 'defer' guarantees that this func() will always be called,
		// even after a panic event.
		defer func() {
			// is there a panic to recover from? Well, in that case ...
			if err := recover(); err != nil {
				w.Header().Set(`Connection`, `close`)
				app.ServerError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// vim: ts=4 sw=4 fdm=indent
