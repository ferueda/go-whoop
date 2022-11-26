go-whoop
=======

[![go-whoop release (latest SemVer)](https://img.shields.io/github/v/release/ferueda/go-whoop?sort=semver)](https://github.com/ferueda/go-whoop/releases)
[![GoDoc](https://godoc.org/github.com/ferueda/go-whoop?status.svg)](http://godoc.org/github.com/ferueda/go-whoop)
[![Test Status](https://github.com/ferueda/go-whoop/workflows/tests/badge.svg)](https://github.com/ferueda/go-whoop/actions?query=workflow%3Atests)
[![Go Report Card](https://goreportcard.com/badge/github.com/ferueda/go-whoop)](https://goreportcard.com/report/github.com/ferueda/go-whoop)

go-whoop is a Go client library for accessing the [WHOOP Platform API](https://developer.whoop.com/api).

## Installation

To install the library, simply

`go get github.com/ferueda/go-whoop`

## Usage
```go
import "github.com/ferueda/go-whoop/whoop"
```
Create a new client, then use the various services on the client to access different parts of the API. For example:
```go
import "github.com/ferueda/go-whoop/whoop"

func main() {
    client := whoop.NewClient(nil)
    ctx := context.Background()

    // list all cycles for the authenticated user
    cycles, err := client.Cycle.ListAll(ctx)
}

```

## Authentication
The client does not handle authentication for you. Instead, you can provide `whoop.NewClient()` with an `http.Client` of your own that can handle authentication for you.

The most common way is using the [OAuth2 package](https://pkg.go.dev/golang.org/x/oauth2).

If you have an OAuth2 access token, you can use it with the OAuth2 package like:

```go
import (
  "golang.org/x/oauth2"
  "github.com/ferueda/go-whoop/whoop"
)

func main() {
    ctx := context.Background()
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: "your_token"},
    )

    client := github.NewClient(oauth2.NewClient(ctx, ts))
    cycles, err := client.Cycle.ListAll(ctx)
}
```

## How to Contribute

* Fork a repository
* Add/Fix something
* Check that tests are passing
* Create PR

Current contributors:

- [Felipe Rueda](https://github.com/ferueda)

## License ##

This library is distributed under the MIT License found in the [LICENSE](./LICENSE)
file.
