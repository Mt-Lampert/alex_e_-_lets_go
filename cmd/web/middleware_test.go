package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MtLampert/alex_e_-_lets_go/internal/assert"
)

func TestSecureHeaders(t *testing.T) {
	// initialize a Response object and a http request stub
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, `/`, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock HTTP handler that we can pass to our secureHeaders()
	// middleware which writes a 200 OK status code and an 'OK.' response body
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `OK.`)
	})

	secureHeaders(next).ServeHTTP(rr, r)

	// get the results of the 'request'
	rs := rr.Result()
	defer rs.Body.Close()

	// test #01
	expected := `default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com`
	assert.Equal(t, rs.Header.Get(`Content-Security-Policy`), expected)
	// test #02
	expected = `origin-when-cross-origin`
	assert.Equal(t, rs.Header.Get(`Referrer-Policy`), expected)
	// test #03
	expected = `nosniff`
	assert.Equal(t, rs.Header.Get(`X-content-Type-Options`), expected)
	// test #04
	expected = `deny`
	assert.Equal(t, rs.Header.Get(`X-Frame-Options`), expected)
	// test #05
	expected = `0`
	assert.Equal(t, rs.Header.Get(`X-XSS-Protection`), expected)

	//------------------------------------------------------------
	// test whether secureHeaders has passed the staff properly to 'next'
	//------------------------------------------------------------

	// Check the Status code
	assert.Equal(t, rs.StatusCode, http.StatusOK)

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
