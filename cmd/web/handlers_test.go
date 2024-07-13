package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MtLampert/alex_e_-_lets_go/internal/assert"
)

func TestPing(t *testing.T) {
	// Initialize an httpTest.Recorder instance
	rr := httptest.NewRecorder()

	// Initialize a stub http Request
	r, err := http.NewRequest(http.MethodGet, `/`, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Call the 'ping' handler function, using 'rr' as writer and 'r' as request.
	ping(rr, r)

	// Get the response object for the Handler Call
	rs := rr.Result()
	// ensure that the body of the result will be garbage collected
	defer rs.Body.Close()

	// Check the Status code
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	//------------------------------------------------------------

	// use IO to read the request body and save it in the 'body' variable as
	// []bytes Slice
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	// remove all leading and trailing whitespace
	bytes.TrimSpace(body)

	assert.Equal(t, string(body), `OK.`)
}

// vim: ts=4 sw=4 fdm=indent
