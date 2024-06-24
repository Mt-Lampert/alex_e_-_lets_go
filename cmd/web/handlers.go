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

// type for saving and validating form data for use in a snippet Template
type SnippetCreateForm struct {
	Title       string
	Content     string
	Expires     string
	FieldErrors map[string]string
}

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

// A handler function to show a "Create Snippet" form
func (app Application) handleNewSnippetForm(w http.ResponseWriter, r *http.Request) {
	data := app.buildTemplateData()
	// will check 'One Month' in the 'Expires' section in the template
	data.Form = SnippetCreateForm{
		Expires: `1 month`,
	}
	app.Render(w, http.StatusOK, `createSnippet.go.html`, data)
}

// A handler function for creating a snippet in the database
func (app Application) handleNewSnippet(w http.ResponseWriter, r *http.Request) {
	// TODO: Validation
	//   1. implement `func validateTitle(rawTitle string) bool {}`
	//   0. implement `func validateContent(rawContent string) bool {}`
	//   0. implement `func validateExpires(rawExpires string) bool {}`
	//   0. Check and 'punish' validation errors

	// ctx := context.Background()
	// Get the form values from the request
	r.ParseForm()
	form := SnippetCreateForm{
		Title:       r.Form.Get("title"),
		Content:     r.Form.Get("content"),
		Expires:     r.Form.Get(`expires`),
		FieldErrors: make(map[string]string, 3),
	}

	//
	// Validate each and every form field
	//
	if !validateTitle(form.Title) {
		form.FieldErrors[`Title`] = `Title entry must be between 4 and 30 characters long.`
	}
	if !validateContent(form.Content) {
		form.FieldErrors[`Content`] = "Content entry must be at least 5 characters long!"
	}

	if !validateExpires(form.Expires) {
		form.FieldErrors[`Expires`] = "Expires entry must be one of the choices below!"
	}
	// evaluate the errors
	if len(form.FieldErrors) > 0 {
		data := app.buildTemplateData()
		data.Form = form
		app.Render(w, http.StatusUnprocessableEntity, `createSnippet.go.html`, data)
		return
	}
	//
	// params2Insert := db.InsertSnippetParams{
	// 	Title:   form.Title,
	// 	Content: form.Content,
	// 	Expires: sql.NullString{Valid: true, String: form.Expires},
	// }
	//
	// feedback, err := db.Qs.InsertSnippet(ctx, params2Insert)
	// if err != nil {
	// 	app.ServerError(w, err)
	// 	return
	// }
	//
	// url := fmt.Sprintf("/snippets/%d", feedback.ID)
	//
	// http.Redirect(w, r, url, http.StatusSeeOther)
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
