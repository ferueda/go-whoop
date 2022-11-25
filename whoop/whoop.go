// Package whoop provides utilties for interfacing
// with the WHOOP API.
package whoop

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	baseURL = "https://api.prod.whoop.com/"

	headerRateLimit     = "X-RateLimit-Limit"
	headerRateRemaining = "X-RateLimit-Remaining"
	headerRateReset     = "X-RateLimit-Reset"
)

type service struct {
	client *Client
}

// Client manages communication with the WHOOP API.
type Client struct {
	http    *http.Client // HTTP client used to communicate with the API.
	baseURL *url.URL     // Base URL for API requests.

	rateLimit Rate // Rate limit for the client as determined by the most recent API call.

	shared service // Reuse a single struct instead of allocating one for each service on the heap.

	// Services used for talking to different parts of the API.
}

// NewClient returns a new WHOOP API client.
// If a nil httpClient is provided, a new http.Client will be used.
//
// To use API methods which require authentication,
// provide an http.Client that will perform the authentication
// for you (such as that provided by the golang.org/x/oauth2 library).
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	baseURL, _ := url.Parse(baseURL)

	c := &Client{http: httpClient, baseURL: baseURL}
	c.shared.client = c
	return c
}

// newRequest creates a new API request with context. If specified,
// the value pointed to by body is JSON encoded and included in the request body.
func (c *Client) newRequest(ctx context.Context, method, url string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(url)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// Response is a WHOOP API response. This wraps the standard http.Response
// returned from the API and provides convenient access to things like
// pagination tokens and rate limits.
type Response struct {
	*http.Response

	// The WHOOP API implements pagination through cursors.
	// This means that a token points directly to the next set of records.
	//
	// For more details visit https://developer.whoop.com/docs/developing/pagination
	NextPageToken string

	Rate Rate

	// token's expiration date
	TokenExpiration time.Time
}

// newResponse creates a new Response for the provided http.Response.
// r must not be nil.
func newResponse(r *http.Response) *Response {
	response := Response{Response: r}
	response.Rate = parseRateLimit(r)
	return &response
}

// TODO
// parseTokenExpiration parses the TokenExpiration related headers.
// func parseTokenExpiration(r *http.Response) Timestamp {
// 	var exp Timestamp
// 	if v := r.Header.Get(headerTokenExpiration); v != "" {
// 		if t, err := time.Parse("2006-01-02 03:04:05 MST", v); err == nil {
// 			exp = Timestamp{t.Local()}
// 		}
// 	}
// 	return exp
// }

// Rate represents the rate limit for the current client.
// After an API key's rate limit has been reached or exceeded,
// the API will respond with a 429 - Too Many Requests HTTP status code.
//
// For more details on how the API handles rate limits
// https://developer.whoop.com/docs/developing/rate-limiting#rate-limit-information
type Rate struct {
	// The number of remaining requests the client can make within the time window.
	// https://developer.whoop.com/docs/developing/rate-limiting#x-ratelimit-remaining
	Remaining int `json:"remaining"`

	// The time at which the current rate limit will reset.
	Reset time.Time `json:"reset"`
}

// parseRateLimit returns the rate limits for the current client.
// It extracts rate limit values from the response headers.
func parseRateLimit(r *http.Response) Rate {
	var rate Rate
	if remaining := r.Header.Get(headerRateRemaining); remaining != "" {
		rate.Remaining, _ = strconv.Atoi(remaining)
	}
	if reset := r.Header.Get(headerRateReset); reset != "" {
		v, _ := strconv.Atoi(reset)
		rate.Reset = now().Add(time.Second * time.Duration(v))
	}
	return rate
}

// do sends the request to the API.
// The response body will be unmarshalled into v,
// or return an error if an API error ocurred.
func (c *Client) do(req *http.Request, v any) (*Response, error) {
	if err := c.checkRateLimit(req); err != nil {
		return &Response{
			Response: err.Response,
			Rate:     err.Rate,
		}, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	response := newResponse(resp)
	defer response.Body.Close()
	c.rateLimit = response.Rate
	return response, json.NewDecoder(response.Body).Decode(v)
}

// get makes a GET request to the given url. The response body will be
// unmarshalled into v.
func (c *Client) get(ctx context.Context, url string, v any) error {
	req, err := c.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, v)
	if err != nil {
		return err
	}
	return nil
}

// checkRateLimit validates if API rate limits have been
// reached or exceeded, for the current client.
// It returns a RateLimitError with a fake response value if
// any rate limit has been reached or nil if not.
//
// Note that it skips making actual network requests
// if rate limits have been reached or exceeded.
func (c *Client) checkRateLimit(req *http.Request) *RateLimitError {
	if !c.rateLimit.Reset.IsZero() && c.rateLimit.Remaining <= 0 && now().Before(c.rateLimit.Reset) {
		// Create a fake response.
		resp := &http.Response{
			Status:     http.StatusText(http.StatusTooManyRequests),
			StatusCode: http.StatusTooManyRequests,
			Request:    req,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
		}
		return &RateLimitError{
			Rate:     c.rateLimit,
			Response: resp,
			Message:  fmt.Sprintf("API rate limit has been reached or exceeded. Please try again after %v.", c.rateLimit.Reset.Format("2006-01-02T15:04:05")),
		}
	}
	return nil
}

// RateLimitError occurs when the API returns 429 Too Many Requests response
// with a rate limit remaining value of 0.
type RateLimitError struct {
	Rate     Rate           // Rate specifies last known rate limit for the client
	Response *http.Response // HTTP response that caused the error
	Message  string         `json:"message"`
}

func (r *RateLimitError) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message)
}

// Is returns whether the provided error equals this error.
func (r *RateLimitError) Is(target error) bool {
	v, ok := target.(*RateLimitError)
	if !ok {
		return false
	}

	return r.Rate == v.Rate &&
		r.Message == v.Message && r.Response != nil && v.Response != nil && r.Response.StatusCode == v.Response.StatusCode
}

// This helper method is useful for testing purposes only.
var now = func() time.Time {
	return time.Now()
}
