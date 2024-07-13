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
	"github.com/mattn/go-sqlite3"
)

// type for saving and validating form data for use in a snippet Template
type SnippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             string `form:"expires"`
	validator.Validator `form:"-"`
}

type SignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type LoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
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
	data := app.buildTemplateData(r)

	data.Snippets = tplSnippets

	app.Render(w, http.StatusOK, `home.go.html`, data)
}

// A handler function to show a "Create Snippet" form
func (app Application) handleNewSnippetForm(w http.ResponseWriter, r *http.Request) {
	data := app.buildTemplateData(r)
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
		data := app.buildTemplateData(r)
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
	app.sessionManager.Put(r.Context(), `Flash`, `New snippet successfully created!`)

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
	flash := app.sessionManager.PopString(r.Context(), `Flash`)

	// create data object
	data := app.buildTemplateData(r)
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

func (app *Application) handleUrlQuery(w http.ResponseWriter, r *http.Request) {
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

func (app *Application) handleUnderConstruction(w http.ResponseWriter, r *http.Request) {
	data := app.buildTemplateData(r)
	data.URL = fmt.Sprintf("GET %s", r.RequestURI)

	app.Render(w, http.StatusOK, `under_construction.go.html`, data)
}

func (app *Application) handleSignupForm(w http.ResponseWriter, r *http.Request) {
	data := app.buildTemplateData(r)
	data.Form = SignupForm{}
	app.Render(w, http.StatusOK, `signup.go.html`, data)
}

func (app *Application) handleSignup(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var form SignupForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	// validate the form fields
	form.CheckField(form.NotBlank(form.Name), `Name`, `'Name' cannot be blank!`)
	form.CheckField(form.NotBlank(form.Email), `Email`, `'Email' cannot be blank!`)
	form.CheckField(form.Matches(form.Email, validator.EmailRegex), `email`, `'Email' must be a valid email address!`)
	form.CheckField(form.NotBlank(form.Password), `Password`, `'Password' cannot be blank!`)
	form.CheckField(form.MinChars(form.Password, 8), `Password`, `'Password' must be at least 8 characters long!`)

	if !form.Valid() {
		data := app.buildTemplateData(r)
		data.Form = form
		app.Render(w, http.StatusUnprocessableEntity, `signup.go.html`, data)
		return
	}

	//
	// insert user data into database
	// ==============================
	//
	// create a Hashed Password from form.Password
	hashedPassword, err := encryptPassword(form.Password)
	if err != nil {
		app.ServerError(w, err)
		return
	}
	// build insertion object
	iup := db.InsertUserParams{
		Name:           form.Name,
		Email:          form.Email,
		HashedPassword: hashedPassword,
	}
	// insert into DB
	_, err = db.Qs.InsertUser(ctx, iup)
	if err != nil {
		// is it an sqlite3 error?
		if sqlite3Err, ok := err.(sqlite3.Error); ok {
			// is it an sqlite3.ErrConstraint error?
			if sqlite3Err.Code == sqlite3.ErrNo(sqlite3.ErrConstraint) {
				// fmt.Println(`    DB Constraint Error!`)
				form.AddFieldError(`Email`, `Sorry, this email address is already taken!`)
				data := app.buildTemplateData(r)
				data.Form = form
				app.Render(w, http.StatusUnprocessableEntity, `signup.go.html`, data)
				return
			}
		}
	}

	// apply the good news
	app.sessionManager.Put(r.Context(), `Flash`, `Your Signup was successful. Please log in!`)

	// redirect to login
	http.Redirect(w, r, `/user/login_form`, http.StatusSeeOther)
}

func (app *Application) handleLoginForm(w http.ResponseWriter, r *http.Request) {
	data := app.buildTemplateData(r)
	data.Form = LoginForm{}
	// data.URL = r.RequestURI
	app.Render(w, http.StatusOK, `login.go.html`, data)
}

func (app *Application) handleLogin(w http.ResponseWriter, r *http.Request) {
	// decode the data from the form in the request into LoginForm
	var form LoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.ClientError(w, http.StatusBadRequest)
		return
	}

	// Validations
	form.CheckField(form.NotBlank(form.Email), `Email`, `Email cannot be blank!`)
	form.CheckField(form.Matches(form.Email, validator.EmailRegex), `Email`, `Please enter a regular email!`)
	form.CheckField(form.NotBlank(form.Password), `Password`, `Password cannot be blank!`)

	if !form.Valid() {
		data := app.buildTemplateData(r)
		data.Form = form
		app.Render(w, http.StatusUnprocessableEntity, `login.go.html`, data)
		return
	}

	// check the credentials
	userId, err := Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			form.AddNonFieldError(`Email or Password is incorrect.`)
			data := app.buildTemplateData(r)
			data.Form = form
			app.Render(w, http.StatusUnprocessableEntity, `login.go.html`, data)
		} else {
			app.ServerError(w, err)
		}
		return
	}

	// renew the session token or create it if it doesn't exist
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Since we made it here, add the ID of the user to the session.
	// Now the user is officially ‘logged in’
	app.sessionManager.Put(r.Context(), `userID`, userId)

	// Redirect the user to the ‘Create Snippet’ page
	http.Redirect(w, r, `/new/snippet`, http.StatusSeeOther)
}

func (app *Application) handleLogout(w http.ResponseWriter, r *http.Request) {
	// First, renew the token!
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.ServerError(w, err)
		return
	}
	// this ‘logs out’ the user!
	app.sessionManager.Remove(r.Context(), `userID`)
	// flash message to inform the user
	app.sessionManager.Put(r.Context(), `Flash`, `You have been successfully logged out!`)
	// redirect to the landing page
	http.Redirect(w, r, `/`, http.StatusSeeOther)
}

// ping for testing
func ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `OK.`)
}

// vim: ts=4 sw=4 fdm=indent
