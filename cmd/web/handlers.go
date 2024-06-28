package main

import (
	"context"

	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MtLampert/alex_e_-_lets_go/internal/db"
	"github.com/MtLampert/alex_e_-_lets_go/internal/validator"
	"github.com/go-chi/chi/v5"
)

// type for saving and validating form data for use in a snippet Template
type SnippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             string `form:"expires"`
	validator.Validator `form:"-"`
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
	ctx := context.Background()
	// simple declaration; serves as target for app.formDecoder.Decode()
	var form SnippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	err = app.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		fmt.Println(`    -> Decode Error`)
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	//
	// Validate each and every form field
	//
	form.CheckField(
		form.WithinRange(4, 20, form.Title),
		`Title`,
		`Title entry must be between 4 and 30 characters long.`,
	)
	form.CheckField(
		form.LongEnough(5, form.Content),
		`Content`,
		"Content entry must be at least 5 characters long!",
	)
	form.CheckField(
		form.ValidExpiration(form.Expires),
		`Expires`,
		"Expires entry must be one of the choices below!",
	)

	// for key, val := range form.FieldErrors {
	// 	fmt.Printf("   '%s': '%s'\n", key, val)
	// }

	// evaluate the errors
	if !form.Valid() {
		// fmt.Println("There are form errors here!")
		data := app.buildTemplateData()
		data.Form = form
		app.Render(w, http.StatusUnprocessableEntity, `createSnippet.go.html`, data)
		return
	}
	//
	params2Insert := db.InsertSnippetParams{
		Title:   form.Title,
		Content: form.Content,
		Expires: sql.NullString{Valid: true, String: form.Expires},
	}

	feedback, err := db.Qs.InsertSnippet(ctx, params2Insert)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// add flash message to the session
	app.sessionManager.Put(r.Context(), `flash`, `New snippet successfully created!`)

	url := fmt.Sprintf("/snippets/%d", feedback.ID)

	http.Redirect(w, r, url, http.StatusSeeOther)
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
			fmt.Println("   Eh? Nothing Found?")
			app.NotFound(w)
		} else {
			app.ServerError(w, err)
		}
		return
	}
	resultTpl := app.ResultRawToTpl(resultRaw)

	// get the flash message (and delete it from the session), if there is one
	flash := app.sessionManager.PopString(r.Context(), `flash`)

	// create data object
	data := app.buildTemplateData()
	data.Snippet = resultTpl
	data.Flash = flash

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
