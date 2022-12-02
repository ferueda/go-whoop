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
client := whoop.NewClient(nil)
ctx := context.Background()

// list all cycles for the authenticated user
cycles, _ := client.Cycle.ListAll(ctx, nil)
```

### Query filters

Some API methods have optional parameters that can be passed to filter results by dates, limit the number of results returned, or provied the token for the next page of results. For example:

```go
client := whoop.NewClient(nil)
ctx := context.Background()

// List all cycle records for the authenticated user with query filters.
params := whoop.RequestParams{
		Start:     time.Now().Add(time.Hour * -48),
		End:       time.Now(),
		Limit:     4,
		NextToken: "abc"}

cycles, err := client.Cycle.ListAll(ctx, &params)
```

### User Service
Get the profile for the authenticated user.
```go
profile, err := client.User.GetProfile(ctx)
```
Get the body measurements for the authenticated user.
```go
bodyMeasurement, err := client.User.GetBodyMeasurement(ctx)
```

### Cycle Service
Get a single physiological cycle record for the specified id.
```go
cycle, err := client.Cycle.GetOne(ctx, 1)
```
List all physiological cycle records for the authenticated user.
```go
cycles, err := client.Cycle.ListAll(ctx, nil)
```

### Sleep Service
Get a single single sleep record for the specified id.
```go
sleep, err := client.Sleep.GetOne(ctx, 1)
```
List all sleep records for the authenticated user.
```go
sleeps, err := client.Sleep.ListAll(ctx, nil)
```

### Recovery Service
Get a single recovery record for the specified cycle id.
```go
recovery, err := client.Recovery.GetOneByCycleId(ctx, 1)
```
List all recovery records for the authenticated user.
```go
recoveries, err := client.Recovery.ListAll(ctx, nil)
```

### Workout Service
Get a single workout activity record for the specified id.
```go
workout, err := client.Workout.GetOne(ctx, 1)
```
List all workout activity records for the authenticated user.
```go
workouts, err := client.Workout.ListAll(ctx, nil)
```

## Authentication
The client does not handle authentication for you. Instead, you can provide `whoop.NewClient()` with an `http.Client` of your own that can handle authentication for you.

The most common way is using the [OAuth2 package](https://pkg.go.dev/golang.org/x/oauth2).

If you have an OAuth2 access token, you can use it with the OAuth2 package like:

```go
import (
    "context"

    "github.com/ferueda/go-whoop/whoop"
    "golang.org/x/oauth2"
)

func main() {
    ctx := context.Background()
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: "your_token"},
    )

    client := whoop.NewClient(oauth2.NewClient(ctx, ts))
    cycles, _ := client.Cycle.ListAll(ctx, nil)
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
