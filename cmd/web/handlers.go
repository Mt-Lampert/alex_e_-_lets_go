package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MtLampert/alex_e_-_lets_go/internal/db"
	"github.com/go-chi/chi/v5"
)

func (app *Application) handleHome(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	rawSnippets, err := db.Qs.GetAllSnippets(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.NotFound(w)
		} else {
			app.ServerError(w, err)
		}
		return
	}
	tplSnippets := app.RawSnippetsToTpl(rawSnippets)

	// create data object
	data := app.buildTemplateData()
	data.Snippets = tplSnippets

	app.Render(w, http.StatusOK, `home.go.html`, data)
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

	idBE := chi.URLParam(r, `id`)
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

	// create data object
	data := app.buildTemplateData()
	data.Snippet = resultTpl

	app.Render(w, http.StatusOK, `view.go.html`, data)
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
