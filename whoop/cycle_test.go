package whoop

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestCycleService_ListAll(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/"+apiVersion+cycleEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testHttpMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `
		{
			"records": [
					{
							"id": 1,
							"user_id": 1,
							"created_at": "2022-11-27T16:34:36.226Z",
							"updated_at": "2022-11-27T16:34:36.226Z",
							"start": "2022-11-27T08:17:45.687Z",
							"end": null,
							"timezone_offset": "-08:00",
							"score_state": "SCORED",
							"score": {
									"strain": 4.0210266,
									"kilojoule": 6350.486,
									"average_heart_rate": 52,
									"max_heart_rate": 115
							}
					},
					{
							"id": 2,
							"user_id": 1,
							"created_at": "2022-11-26T15:17:19.570Z",
							"updated_at": "2022-11-27T16:34:41.484Z",
							"start": "2022-11-26T06:49:28.470Z",
							"end": "2022-11-27T08:17:45.687Z",
							"timezone_offset": "-08:00",
							"score_state": "SCORED",
							"score": {
									"strain": 6.1106668,
									"kilojoule": 9490.08,
									"average_heart_rate": 55,
									"max_heart_rate": 121
							}
					}
			],
			"next_token": null
	}`)
	})

	ctx := context.Background()
	resp, err := client.Cycle.ListAll(ctx, nil)

	if err != nil {
		t.Fatalf("Cycle.ListAll(): expected nil error, got %#v", err)
	}
	if len(resp.Records) != 2 {
		t.Errorf("Cycle.ListAll(): expected 2 records, got %v", len(resp.Records))
	}
	if resp.NextToken != nil {
		t.Errorf("Cycle.ListAll(): expected next_token nil, got %v", resp.NextToken)
	}
}

func TestCycleService_ListAll_with_params(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	date := now()
	mux.HandleFunc("/"+apiVersion+cycleEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testHttpMethod(t, r, http.MethodGet)
		limit := r.URL.Query().Get("limit")
		nextToken := r.URL.Query().Get("nextToken")
		start := r.URL.Query().Get("start")
		if limit != "1" {
			t.Errorf("Cycle.ListAll(): expected limit == 1, got %v", limit)
		}
		if nextToken != "test_token" {
			t.Errorf("Cycle.ListAll(): expected nextToken == test_token, got %v", nextToken)
		}
		if start != date.Format(time.RFC3339) {
			t.Errorf("Cycle.ListAll(): expected start == %v, got %v", date.Format(time.RFC3339), start)
		}
		fmt.Fprint(w, `
		{
			"records": [
					{
							"id": 1,
							"user_id": 1,
							"created_at": "2022-11-27T16:34:36.226Z",
							"updated_at": "2022-11-27T16:34:36.226Z",
							"start": "2022-11-27T08:17:45.687Z",
							"end": null,
							"timezone_offset": "-08:00",
							"score_state": "SCORED",
							"score": {
									"strain": 4.0210266,
									"kilojoule": 6350.486,
									"average_heart_rate": 52,
									"max_heart_rate": 115
							}
					}
			],
			"next_token": "test_token"
	}`)
	})

	ctx := context.Background()
	params := RequestParams{Limit: 1, NextToken: "test_token", Start: date}
	resp, err := client.Cycle.ListAll(ctx, &params)

	if err != nil {
		t.Fatalf("Cycle.ListAll(): expected nil error, got %#v", err)
	}
	if len(resp.Records) != 1 {
		t.Errorf("Cycle.ListAll(): expected 2 records, got %v", len(resp.Records))
	}
	if resp.NextToken == nil {
		t.Errorf("Cycle.ListAll(): expected next_token == test_token, got %v", resp.NextToken)
	}
}

func TestCycleService_GetOne(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/"+apiVersion+cycleEndpoint+"/", func(w http.ResponseWriter, r *http.Request) {
		testHttpMethod(t, r, http.MethodGet)
		urlPath := strings.Split(r.URL.Path, "/")
		if len(urlPath) != 4 {
			t.Errorf("Cycle.GetOne(): expected url path of length 4, got %v", len(urlPath))
		}
		if urlPath[len(urlPath)-1] != "1" {
			t.Errorf("Cycle.GetOne(): expected id 1, got %v", urlPath[len(urlPath)-1])
		}
		fmt.Fprint(w, `
			{
					"id": 1,
					"user_id": 1,
					"created_at": "2022-11-27T16:34:36.226Z",
					"updated_at": "2022-11-27T16:34:36.226Z",
					"start": "2022-11-27T08:17:45.687Z",
					"end": null,
					"timezone_offset": "-08:00",
					"score_state": "SCORED",
					"score": {
							"strain": 4.0210266,
							"kilojoule": 6350.486,
							"average_heart_rate": 52,
							"max_heart_rate": 115
					}
			}
	`)
	})

	ctx := context.Background()
	resp, err := client.Cycle.GetOne(ctx, "1")

	if err != nil {
		t.Fatalf("Cycle.GetOne(): expected nil error, got %#v", err)
	}
	if resp.ID != 1 {
		t.Errorf("Cycle.GetOne(): expected id 1, got %v", resp.ID)
	}
}
