package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
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

func (app *Application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, `/user/login_form`, http.StatusSeeOther)
		}
		// If we made it here, set the 'Cache-Control: no store' header
		// so that pages that require authentication are not stored in the
		// user's browser cache (or other caches they employ)
		w.Header().Add(`Cache-Control`, `no store`)
		next.ServeHTTP(w, r)
	})
}

func noSurf(next http.Handler) http.Handler {
	crsfHandler := nosurf.New(next)
	crsfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     `/`,
		Secure:   true,
	})
	return crsfHandler
}

func (app *Application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// do we already have an authenticated userID in the session?
		id := app.sessionManager.GetInt64(r.Context(), `userID`)
		// if not found
		if id == 0 {
			// we needn't check further, so ...
			next.ServeHTTP(w, r)
			return
		}

		// Ok, we have to check!
		exists, err := app.userExists(id)
		if err != nil {
			app.ServerError(w, err)
			return
		}

		// Hey! That user really exists!
		if exists {
			// define an updated context instance
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

// vim: ts=4 sw=4 fdm=indent
