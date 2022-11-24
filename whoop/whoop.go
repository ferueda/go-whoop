// Package whoop provides utilties for interfacing
// with the WHOOP API.
package whoop

import (
	"net/http"
	"net/url"
	"time"
)

const (
	baseURL = "https://api.prod.whoop.com/"
)

type service struct {
	client *Client
}

// Client manages communication with the WHOOP API.
type Client struct {
	http    *http.Client // HTTP client used to communicate with the API.
	baseURL *url.URL     // Base URL for API requests

	shared service // Reuse a single struct instead of allocating one for each service on the heap.

	// Services used for talking to different parts of the API.
}

// NewClient returns a new WHOOP API client.
// If a nil httpClient is provided, a new http.Client will be used.
//
// To use API methods which require authentication, you must call
// the Client.Authenticate method with valid credentials.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: time.Minute}
	}
	baseURL, _ := url.Parse(baseURL)

	c := &Client{http: httpClient, baseURL: baseURL}

	c.shared.client = c
	return c
}
