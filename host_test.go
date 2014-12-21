package pinger

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func ExampleHost() {
	h := Host{
		Name:               "test host",
		Url:                "http://example.com/_ping",
		ExpectedStatusCode: 200,
		ExpectedBody:       "PONG",
	}

	if status, body, err := h.Ping(http.DefaultClient); err != nil {
		log.Print(status, string(body), err)
	}
}

func TestPing(t *testing.T) {
	// test invalid url
	h := Host{
		Name: "test",
		Url:  "http://example",
	}

	if status, body, err := h.Ping(http.DefaultClient); status != 0 || body != nil || err == nil {
		t.Fail()
	}

	// test status code 200
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {}))

	h.Url = ts.URL
	if status, body, err := h.Ping(http.DefaultClient); status != 200 ||
		string(body) != "" || err != nil {

		t.Fail()
	}

	// test status code 200 with body
	ts = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "test body")
		}))

	h.Url = ts.URL
	if status, body, err := h.Ping(http.DefaultClient); status != 200 ||
		string(body) != "test body" || err != nil {

		t.Fail()
	}

	// test bad status code
	ts = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
		}))

	h.Url = ts.URL
	if status, body, err := h.Ping(http.DefaultClient); status != 500 ||
		string(body) != "" || err != ErrBadStatusCode {

		t.Fail()
	}

	// test status code expectation mismatch
	ts = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {}))

	h.Url = ts.URL
	h.ExpectedStatusCode = 500

	if status, body, err := h.Ping(http.DefaultClient); status != 200 ||
		string(body) != "" || err != ErrStatusCodeMismatch {

		t.Fail()
	}

	// test body expectation mismatch
	ts = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "new body")
		}))

	h.Url = ts.URL
	h.ExpectedStatusCode = 0
	h.ExpectedBody = "old body"

	if status, body, err := h.Ping(http.DefaultClient); status != 200 ||
		string(body) != "new body" || err != ErrBodyMismatch {

		t.Fail()
	}
}
