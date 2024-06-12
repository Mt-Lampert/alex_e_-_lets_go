package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/MtLampert/alex_e_-_lets_go/internal/db"
)

type TplSnippet struct {
	ID      string
	Title   string
	Content string
	Created string
	Expires string
}

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

// A convenient wrapper around 'clientError()' which sends a "404 Not Found"
// error to the frontend; it will be needed very often!
func (app *Application) NotFound(w http.ResponseWriter) {
	app.ClientError(w, http.StatusNotFound)
}

// converts a single DB snippet into a snippet for use in templates
func (app *Application) ResultRawToTpl(r db.GetSnippetRow) TplSnippet {
	id := strconv.Itoa(int(r.ID))
	var createdTpl string

	if r.Created.Valid {
		createdTpl = r.Created.Time.Format("2006-01-02 03:04:05")
	}

	myExpiresTpl := fmt.Sprintf("%v", r.Ends)

	return TplSnippet{
		ID:      id,
		Title:   r.Title,
		Content: r.Content,
		Created: createdTpl,
		Expires: myExpiresTpl,
	}

}

// converts a slice of raw DB snippets into snippets for use in templates
func (app *Application) RawSnippetsToTpl(rs []db.GetAllSnippetsRow) []TplSnippet {
	var createdTpl string
	lenRS := cap(rs)
	var tsp = make([]TplSnippet, lenRS)

	for i, r := range rs {
		id := strconv.Itoa(int(r.ID))
		if r.Created.Valid {
			createdTpl = r.Created.Time.Format("2006-01-02 03:04:05")
		}
		myExpiresTpl := fmt.Sprintf("%v", r.Ends)
		tsp[i] = TplSnippet{
			ID:      id,
			Title:   r.Title,
			Content: r.Content,
			Created: createdTpl,
			Expires: myExpiresTpl,
		}
	}
	return tsp
}

// renders templates and sends them to the frontend
func (app *Application) Render(
	w http.ResponseWriter,
	status int,
	page string,
	data *templateData,
) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template '%s' does not exist", page)
		app.ServerError(w, err)
		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, `base`, data)
	if err != nil {
		app.ServerError(w, err)
	}
}

// vim: ts=4 sw=4 fdm=indent
