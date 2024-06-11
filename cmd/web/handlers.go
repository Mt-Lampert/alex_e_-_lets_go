package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/MtLampert/alex_e_-_lets_go/internal/db"
)

func (app *Application) handleHome(w http.ResponseWriter, r *http.Request) {
	templates := []string{
		"./ui/html/base.go.html",
		"./ui/html/pages/home.go.html",
		"./ui/html/partials/nav.go.html",
	}

	// See Journal, 2024-06-03 07:43 for documentation
	ts, err := template.ParseFiles(templates...)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// See Journal, 2024-06-03 07:43 for documentation
	err = ts.ExecuteTemplate(w, `base`, nil)
	if err != nil {
		app.ServerError(w, err)
		return
	}
}

// Add a handler function for creating a snippet.
func (app Application) handleNewSnippet(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	output := "Inserted new Snippet:\n"
	params2Insert := db.InsertSnippetParams{
		Title:   `So ein Dummy`,
		Content: `Ich bin ja so ein Dummy!`,
		Expires: sql.NullString{Valid: true, String: `30 days`},
	}

	feedback, err := db.Qs.InsertSnippet(ctx, params2Insert)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.InfoLog.Println("Inserted new entry!")

	output += fmt.Sprintf("    id: %d\n", feedback.ID)
	output += fmt.Sprintf("    title: '%s'\n\n", feedback.Title)

	fmt.Fprint(w, output)
}

// Add a handler function for viewing a specific snippet
func (app Application) handleSingleSnippetView(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	idBE := r.PathValue(`id`)
	idDB, err := strconv.ParseInt(idBE, 10, 64)
	if err != nil {
		app.NotFound(w)
		return
	}

	resultRaw, err := db.Qs.GetSnippet(ctx, idDB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.NotFound(w)
		} else {
			app.ServerError(w, err)
		}
		return
	}

	resultTpl := app.ResultRawToTpl(resultRaw)

	data := &templateData{
		Snippet: resultTpl,
	}

	// Initialize a slice containing the paths to the 'view.go.html' file
	// plus the base layout and navigation partial that we made earlier.
	myTemplates := []string{
		"./ui/html/base.go.html",
		"./ui/html/partials/nav.go.html",
		"./ui/html/pages/view.go.html",
	}

	// Parse the templates
	ts, err := template.ParseFiles(myTemplates...)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Now execute them. Notice how we pass in the snippet data and the final
	// parameter
	if err = ts.ExecuteTemplate(w, "base", data); err != nil {
		app.ServerError(w, err)
	}
}

func (app *Application) handleSnippetList(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	output := "We found snippets. Hallelujah!\n"

	resultsRaw, err := db.Qs.GetAllSnippets(ctx)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	app.InfoLog.Println(`snippets found!`)
	for _, sn := range resultsRaw {
		output += fmt.Sprintf("       id: '%d'\n", sn.ID)
		output += fmt.Sprintf("    title: '%s'\n\n", sn.Title)
	}

	fmt.Fprint(w, output)
}

func (app Application) handleUrlQuery(w http.ResponseWriter, r *http.Request) {
	// get the id
	rawId := r.URL.Query().Get(`id`)

	// validate the id:
	//    - it must be numerical
	//    - it must be greater than 0
	id, err := strconv.Atoi(rawId)
	if err != nil || id <= 0 {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, `Seems you are looking for something with id '%s'`, rawId)
}

// vim: ts=4 sw=4 fdm=indent
