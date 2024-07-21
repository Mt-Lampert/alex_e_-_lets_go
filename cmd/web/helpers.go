package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/MtLampert/alex_e_-_lets_go/internal/db"
	"github.com/MtLampert/alex_e_-_lets_go/internal/tpl"
	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
	"golang.org/x/crypto/bcrypt"
)

// type TplSnippet struct {
// 	ID      string
// 	Title   string
// 	Content string
// 	Created string
// 	Expires string
// }

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
func (app *Application) ResultRawToTpl(r db.GetSnippetRow) tpl.DbSnippet {
	id := strconv.Itoa(int(r.ID))
	var createdTpl string

	if r.Created.Valid {
		createdTpl = r.Created.Time.Format("2006-01-02 03:04:05")
	}

	myExpiresTpl := fmt.Sprintf("%v", r.Ends)
	contentTpl := strings.ReplaceAll(r.Content, "\\n", "\n")

	return tpl.DbSnippet{
		ID:      id,
		Title:   r.Title,
		Content: contentTpl,
		Created: createdTpl,
		Expires: myExpiresTpl,
	}
}

// converts a slice of raw DB snippets into snippets for use in templates
func (app *Application) RawSnippetsToTpl(rs []db.Snippet) []tpl.DbSnippet {
	var createdTpl string
	var tsp = make([]tpl.DbSnippet, 0)

	for _, r := range rs {
		if snippetExpired(r.Created.Time, r.Expires.String) {
			continue
		}
		id := strconv.Itoa(int(r.ID))
		if r.Created.Valid {
			createdTpl = r.Created.Time.Format("2006-01-02 03:04:05")
			// myExpiresTpl := fmt.Sprintf("%v", r.Ends)
			tsp = append(tsp, tpl.DbSnippet{
				ID:      id,
				Title:   r.Title,
				Content: r.Content,
				Created: createdTpl,
				Expires: r.Expires.String,
			})
		}
	}
	return tsp
}

// renders templates and sends them to the frontend
func (app *Application) Render(
	w http.ResponseWriter,
	status int,
	page string,
	data *tpl.TemplateData,
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

// factory helper to build a templateData instance
func (app *Application) buildTemplateData(r *http.Request) *tpl.TemplateData {
	return &tpl.TemplateData{
		CSRFToken:       nosurf.Token(r),
		CurrentYear:     strconv.FormatInt(int64(time.Now().Year()), 10),
		IsAuthenticated: app.isAuthenticated(r),
	}
	// erg. Flash:           app.sessionManager.PopString(r.Context(), `flash`),
}

// checks if 'title' form value is valid
func validateTitle(rawTitle string) bool {
	longEnough := utf8.RuneCountInString(rawTitle) >= 4
	shortEnough := utf8.RuneCountInString(rawTitle) <= 30

	return longEnough && shortEnough
}

// checks if 'content' form value is valid
func validateContent(rawContent string) bool {
	return utf8.RuneCountInString(rawContent) >= 5
}

// checks if 'expires' form value is valid
func validateExpires(rawExpires string) bool {
	return rawExpires == `1 day` || rawExpires == `1 week` || rawExpires == `1 month` || rawExpires == `1 year`
}

// catches a `form.invalidDecoderError` that is thrown when an invalid 'dest'
// pointer is passed into formDecoder.Decode() as destination object. Throws
// a panic in this case; otherwise returns any other error or even nil.
func (app *Application) decodePostForm(r *http.Request, dest any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dest, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		// if error happens to be an invalidDecoderError (at the end of the day)
		if errors.As(err, &invalidDecoderError) {
			panic(invalidDecoderError)
		}
		// implicit else
		return err
	}
	return nil
}

// evaluates if a snippet is expired or not
func snippetExpired(created time.Time, expires string) bool {
	now := time.Now()
	timeoutMap := map[string]int{
		`1 day`:   1,
		`1 week`:  7,
		`1 month`: 30,
		`1 year`:  365,
	}

	return created.AddDate(0, 0, timeoutMap[expires]).Before(now)
}

// provides a database connection for session management
func sessionDB() *sql.DB {
	sDB, err := sql.Open(`sqlite3`, `sessions.db`)
	if err != nil {
		log.Fatal(err)
	}
	return sDB
}

// encrypts user password into bcrypt hash
func encryptPassword(pw string) (string, error) {
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(pw), 12)
	if err != nil {
		return "", err
	}
	return string(hashedPW), nil
}

func Authenticate(email, password string) (int64, error) {
	ctx := context.Background()

	userRow, err := db.Qs.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(userRow.HashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return userRow.ID, nil
}

// returns true if user is authenticated, false if not
func (app *Application) isAuthenticated(r *http.Request) bool {
	// return app.sessionManager.Exists(r.Context(), `userID`)
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}

// returns true if a user with an id exists in the database
func (app *Application) userExists(id int64) (bool, error) {
	dbUserCount, err := db.Qs.CheckUser(context.Background(), id)
	if err != nil {
		return false, err
	}

	return dbUserCount == 1, nil
}

// vim: ts=4 sw=4 fdm=indent
