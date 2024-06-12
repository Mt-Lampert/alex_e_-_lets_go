package main

import (
	"html/template"
	"path/filepath"
)

// a convenient wrapper for template data (usually from the DB)
type templateData struct {
	Snippet  TplSnippet
	Snippets []TplSnippet
}

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

		// build a bundle of required templates for a file
		files := []string{
			`./ui/html/base.go.html`,
			`./ui/html/partials/nav.go.html`,
			page,
		}

		// build a corresponding bundle of 'compiled' raw templates
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		// add the template bundle to the cache with the base filename as key,
		// eg 'home.go.html'
		cache[name] = ts
	}

	return cache, nil
}

// vim: ts=4 sw=4 fdm=indent
