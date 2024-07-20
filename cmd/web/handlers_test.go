package main

/*
import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MtLampert/alex_e_-_lets_go/internal/assert"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
)

func TestPing(t *testing.T) {
	sessDB := sessionDB()

	// initializing the scs session manager
	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(sessDB)
	sessionManager.Lifetime = time.Hour * 12
	app := &Application{
		ErrLog:         log.New(io.Discard, "", 0),
		InfoLog:        log.New(io.Discard, "", 0),
		sessionManager: sessionManager,
	}
	// Initialize an httpTest.Recorder instance
	rr := httptest.NewRecorder()

	// Initialize a stub http Request
	r, err := http.NewRequest(http.MethodGet, `/`, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Call the 'ping' handler function, using 'rr' as writer and 'r' as request.
	app.sessionManager.LoadAndSave(app.handlePing(rr, r))

	// Get the response object for the Handler Call
	rs := rr.Result()
	// ensure that the body of the result will be garbage collected
	defer rs.Body.Close()

	// Check the Status code
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	//------------------------------------------------------------

	// use IO to read the request body and save it in the 'body' variable as
	// []bytes Slice
	// body, err := io.ReadAll(rs.Body)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// // remove all leading and trailing whitespace
	// bytes.TrimSpace(body)
	//
	// assert.Equal(t, string(body), `OK.`)
}
*/
// vim: ts=4 sw=4 fdm=indent
