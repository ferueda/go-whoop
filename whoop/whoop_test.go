package whoop

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	if got, want := c.baseURL.String(), baseURL; got != want {
		t.Errorf("NewClient(): baseURL is %v, want %v", got, want)
	}

	c2 := NewClient(nil)
	if c.http == c2.http {
		t.Error("NewClient(): returned same http client, but they should be different")
	}
}

func TestNewRequest(t *testing.T) {
	c := NewClient(nil)

	url := baseURL + apiVersion + "/test"
	req, _ := c.newRequest(context.Background(), http.MethodGet, "/test", nil)

	if got, want := req.URL.String(), url; got != want {
		t.Errorf("NewRequest(%q): URL is %v, want %v", url, got, want)
	}

	if got := req.Body; got != nil {
		t.Errorf("NewRequest(%q): body is %v, want %v", url, got, nil)
	}

	if got, want := req.Header.Get("Accept"), "application/json"; got != want {
		t.Errorf("NewRequest(%q): Accept header is %v, want %v", url, got, want)
	}

	if got, want := req.Header.Get("Content-Type"), ""; got != want {
		t.Errorf("NewRequest(%q): Content-Type header is %v, want %v", url, got, want)
	}
}

func TestParseRateLimit(t *testing.T) {
	date := now()
	now = func() time.Time {
		return date
	}

	testCases := []struct {
		remaining int
		reset     int
	}{
		{60, 30},
		{0, 0},
	}

	for _, test := range testCases {
		res := http.Response{
			Header: http.Header{},
		}

		res.Header.Set(headerRateRemaining, strconv.Itoa(test.remaining))
		res.Header.Set(headerRateReset, strconv.Itoa(test.reset))

		rate := parseRateLimit(&res)

		if got, want := rate.Remaining, test.remaining; got != want {
			t.Errorf("parseRateLimit(): rate.Remaining is %v, want %v", got, want)
		}
		if got, want := rate.Reset, now().Add(time.Duration(test.reset)*time.Second); got != want {
			t.Errorf("parseRateLimit(): rate.Reset is %v, want %v", got, want)
		}
	}
}

func TestNewResponse(t *testing.T) {
	date := now()
	now = func() time.Time {
		return date
	}
	res := http.Response{
		Header: http.Header{},
	}
	res.Header.Set(headerRateRemaining, "60")
	res.Header.Set(headerRateReset, "30")
	response := newResponse(&res)
	if got, want := response.Response, &res; got != want {
		t.Errorf("newResponse(): response.Response is %v, want %v", got, want)
	}
	if got, want := response.NextPageToken, ""; got != want {
		t.Errorf("newResponse(): rate.NextPageToken is %v, want %v", got, want)
	}
	if got, want := response.Rate.Remaining, 60; got != want {
		t.Errorf("newResponse(): rate.Remaining is %v, want %v", got, want)
	}
	if got, want := response.Rate.Reset, now().Add(time.Duration(30)*time.Second); got != want {
		t.Errorf("newResponse(): rate.Reset is %v, want %v", got, want)
	}
}

func TestCheckResponse(t *testing.T) {
	testCases := []struct {
		statusCode int
		body       string
		want       error
	}{
		{200, "", nil},
		{100, "test error", &Error{Code: 100, Message: "test error"}},
		{300, "test error", &Error{Code: 300, Message: "test error"}},
		{400, "test error", &Error{Code: 400, Message: "test error"}},
		{500, "test error", &Error{Code: 500, Message: "test error"}},
	}

	for _, test := range testCases {
		res := http.Response{
			Header:     http.Header{},
			Body:       io.NopCloser(strings.NewReader(test.body)),
			StatusCode: test.statusCode,
		}
		res.Header.Set(headerRateRemaining, "60")
		res.Header.Set(headerRateReset, "30")

		got := checkResponse(&res)
		if test.statusCode != 200 && got.Error() != test.want.Error() {
			t.Errorf("checkResponse(): got %v, want %v", got, test.want)
		}
		if test.statusCode >= 200 && test.statusCode <= 299 && got != test.want {
			t.Errorf("checkResponse(): got %v, want %v", got, test.want)
		}
	}
}

func TestCheckResponse_TooManyRequests(t *testing.T) {
	res := http.Response{
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader("too many requests error")),
		StatusCode: http.StatusTooManyRequests,
	}
	res.Header.Set(headerRateRemaining, "60")
	res.Header.Set(headerRateReset, "30")

	got := checkResponse(&res)
	if got == nil {
		t.Errorf("checkResponse(): got nil error, want 429 error")
	}
	if err, ok := got.(*RateLimitError); !ok {
		t.Errorf("checkResponse(): expected RateLimitError error; got %#v.", err)
	}
}

func TestAddParams(t *testing.T) {
	start := time.Date(2022, 1, 1, 12, 30, 20, 0, time.UTC)
	end := time.Date(2022, 5, 20, 6, 0, 10, 10, time.UTC)
	testCases := []struct {
		url    string
		params *RequestParams
		want   string
	}{
		{"/test", nil, "/test"},
		{"/test", &RequestParams{}, "/test"},
		{"/test", &RequestParams{Limit: 5}, "/test?limit=5"},
		{"/test", &RequestParams{Start: start}, "/test?start=2022-01-01T12%3A30%3A20Z"},
		{"/test", &RequestParams{End: end}, "/test?end=2022-05-20T06%3A00%3A10Z"},
		{"/test", &RequestParams{NextToken: "test_token"}, "/test?nextToken=test_token"},
		{"/test", &RequestParams{Start: start, End: end}, "/test?end=2022-05-20T06%3A00%3A10Z&start=2022-01-01T12%3A30%3A20Z"},
		{"/test", &RequestParams{Limit: 1, NextToken: "test_token"}, "/test?limit=1&nextToken=test_token"},
		{"/test", &RequestParams{Start: start, End: end, Limit: 30, NextToken: "test_token"}, "/test?end=2022-05-20T06%3A00%3A10Z&limit=30&nextToken=test_token&start=2022-01-01T12%3A30%3A20Z"},
	}

	for _, test := range testCases {
		got, _ := addParams(test.url, test.params)

		if got != test.want {
			t.Errorf("addParams(): got %v, want %v", got, test.want)
		}
	}
}

func TestCheckRateLimit(t *testing.T) {
	url, _ := url.Parse("/test")
	date := now()
	now = func() time.Time {
		return date
	}
	c := NewClient(nil)
	testCases := []struct {
		remaining int
		reset     time.Time
		isError   bool
	}{
		{60, now(), false},
		{0, now().Add(time.Minute * 30), true},
		{0, now(), false},
	}

	for _, test := range testCases {
		c.rateLimit.Remaining = test.remaining
		c.rateLimit.Reset = test.reset
		req := http.Request{
			Method: http.MethodGet,
			URL:    url,
			Header: http.Header{},
		}
		got := c.checkRateLimit(&req)

		if !test.isError && got != nil {
			t.Errorf("checkRateLimit(): expected nil; got %#v", got)
		}
		if test.isError && got == nil {
			t.Errorf("checkRateLimit(): expected RateLimitError error; got nil")
		}
	}
}

func TestDo(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	type data struct {
		Test string `json:"test"`
	}

	mux.HandleFunc("/"+apiVersion+"/", func(w http.ResponseWriter, r *http.Request) {
		testHttpMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"test":"test"}`)
	})

	ctx := context.Background()
	req, _ := client.newRequest(ctx, http.MethodGet, "/", nil)
	body := data{}
	err := client.do(req, &body)

	if err != nil {
		t.Errorf("do(): got unexpected error %#v", err)
	}
	if want := (data{"test"}); body != want {
		t.Errorf("do(): response body = %v, want %v", body, want)
	}
}

func TestDo_rateLimit(t *testing.T) {
	date := now()
	now = func() time.Time {
		return date
	}

	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/"+apiVersion+"/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(headerRateRemaining, "60")
		w.Header().Set(headerRateReset, "30")
		fmt.Fprint(w, `{"test":"test"}`)
	})

	ctx := context.Background()
	req, _ := client.newRequest(ctx, http.MethodGet, "/", nil)
	got := client.do(req, &struct {
		Test string `json:"test"`
	}{})

	if got != nil {
		t.Fatalf("do(): got error %#v, expected nil", got)
	}
	if client.rateLimit.Remaining != 60 {
		t.Errorf("do(): expected rateLimit.Remaining 60; got %v.", client.rateLimit.Remaining)
	}
	if client.rateLimit.Reset != now().Add(time.Second*30) {
		t.Errorf("do(): expected rateLimit.Reset %v; got %v.", now().Add(time.Second*30), client.rateLimit.Reset)
	}
}
func TestDo_rateLimit_error(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/"+apiVersion+"/", func(w http.ResponseWriter, r *http.Request) {
		testHttpMethod(t, r, http.MethodGet)
		t.Errorf("do(): should not make the network request")
	})

	ctx := context.Background()
	req, _ := client.newRequest(ctx, http.MethodGet, "/", nil)
	client.rateLimit.Remaining = 0
	client.rateLimit.Reset = now().Add(time.Hour * 5)
	got := client.do(req, &struct{}{})

	if got == nil {
		t.Fatal("do(): got nil error, expected RateLimitError")
	}
	if err, ok := got.(*RateLimitError); !ok {
		t.Errorf("do(): expected RateLimitError error; got %#v.", err)
	}
}

func TestDo_BadRequest(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/"+apiVersion+"/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	ctx := context.Background()
	req, _ := client.newRequest(ctx, http.MethodGet, "/", nil)
	got := client.do(req, &struct{}{})

	if got == nil {
		t.Fatal("do(): got nil error, expected HTTP 400 error")
	}
	if err, ok := got.(*Error); !ok || err.Code != http.StatusBadRequest {
		t.Errorf("do(): expected HTTP 400 error; got %#v.", got)
	}
}

func testHttpMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func setup() (*Client, *http.ServeMux, string, func()) {
	handler := http.NewServeMux()
	server := httptest.NewServer(handler)
	client := NewClient(nil)
	url, _ := url.Parse(server.URL)
	client.baseURL = url
	return client, handler, server.URL, server.Close
}
