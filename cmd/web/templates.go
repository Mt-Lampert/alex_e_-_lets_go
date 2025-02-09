package main

import (
	"html/template"
	"path/filepath"
	"time"
)

// a convenient wrapper for template data (usually from the DB)
// type templateData struct {
// 	CSRFToken       string
// 	CurrentYear     int
// 	Flash           string
// 	Form            any
// 	IsAuthenticated bool
// 	Snippet         TplSnippet
// 	Snippets        []TplSnippet
// 	URL             string
// }

func buildTemplateCache() (map[string]*template.Template, error) {
	// initialize a template cache
	cache := map[string]*template.Template{}

	// collect all the 'pages' template files
	pages, err := filepath.Glob(`./ui/html/pages/*.go.html`)
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// extract the basename of 'page'
		name := filepath.Base(page)

		// parse the base template file into a template set
		ts, err := template.ParseFiles(`./ui/html/base.go.html`)
		if err != nil {
			return nil, err
		}

		// parse all possible partials files into this very template set to add them there
		// notice how the 'old' ts on the right side creates a new enhanced ts
		// on the left hand of '='
		ts, err = ts.ParseGlob(`./ui/html/partials/*.go.html`)
		if err != nil {
			return nil, err
		}

		// parse the current page file into this very template set to add it there
		// notice how the 'old' ts on the right side creates a new enhanced ts
		// on the left hand of '='
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// add the template bundle to the cache with the base filename as key,
		// eg 'home.go.html'
		cache[name] = ts
	}

	return cache, nil
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	// convert t to UTC before formatting it
	return t.UTC().Format(`2006-01-02 03:04`)
}

// vim: ts=4 sw=4 fdm=indent
