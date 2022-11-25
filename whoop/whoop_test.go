package whoop

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	if got, want := c.baseURL.String(), baseURL; got != want {
		t.Errorf("NewClient baseURL is %v, want %v", got, want)
	}

	c2 := NewClient(nil)
	if c.http == c2.http {
		t.Error("NewClient returned same http client, but they should be different")
	}
}

func TestNewRequest(t *testing.T) {
	c := NewClient(nil)

	url := baseURL + "developer/" + apiVersion + "/test"
	req, _ := c.newRequest(context.Background(), http.MethodGet, "/test", nil)

	if got, want := req.URL.String(), url; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", url, got, want)
	}

	if got := req.Body; got != nil {
		t.Errorf("NewRequest(%q) body is %v, want %v", url, got, nil)
	}

	if got, want := req.Header.Get("Accept"), "application/json"; got != want {
		t.Errorf("NewRequest(%q) Accept header is %v, want %v", url, got, want)
	}

	if got, want := req.Header.Get("Content-Type"), ""; got != want {
		t.Errorf("NewRequest(%q) Content-Type header is %v, want %v", url, got, want)
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
			t.Errorf("rate.Remaining is %v, want %v", got, want)
		}
		if got, want := rate.Reset, now().Add(time.Duration(test.reset)*time.Second); got != want {
			t.Errorf("rate.Reset is %v, want %v", got, want)
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
		t.Errorf("response.Response is %v, want %v", got, want)
	}
	if got, want := response.NextPageToken, ""; got != want {
		t.Errorf("rate.NextPageToken is %v, want %v", got, want)
	}
	if got, want := response.Rate.Remaining, 60; got != want {
		t.Errorf("rate.Remaining is %v, want %v", got, want)
	}
	if got, want := response.Rate.Reset, now().Add(time.Duration(30)*time.Second); got != want {
		t.Errorf("rate.Reset is %v, want %v", got, want)
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
			t.Errorf("checkResponse is %v, want %v", got, test.want)
		}
		if test.statusCode >= 200 && test.statusCode <= 299 && got != test.want {
			t.Errorf("checkResponse is %v, want %v", got, test.want)
		}
	}
}

func setup() *Client {
	client := NewClient(nil)
	url, _ := url.Parse("http://test.test/")
	client.baseURL = url
	return client
}
