package whoop

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestRecoveryService_ListAll(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/"+apiVersion+recoveryEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testHttpMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `
		{
			"records": [
					{
						"cycle_id": 1,
						"sleep_id": 1,
						"user_id": 1,
						"created_at": "2022-11-28T13:29:05.961Z",
						"updated_at": "2022-11-28T13:29:12.588Z",
						"score_state": "SCORED",
						"score": {
								"user_calibrating": false,
								"recovery_score": 69.0,
								"resting_heart_rate": 46.0,
								"hrv_rmssd_milli": 89.81702,
								"spo2_percentage": 93.111115,
								"skin_temp_celsius": 33.4
						}
					},
					{
						"cycle_id": 2,
						"sleep_id": 2,
						"user_id": 1,
						"created_at": "2022-11-27T16:34:36.226Z",
						"updated_at": "2022-11-27T16:47:13.425Z",
						"score_state": "SCORED",
						"score": {
								"user_calibrating": true,
								"recovery_score": 97.0,
								"resting_heart_rate": 43.0,
								"hrv_rmssd_milli": 107.67454,
								"spo2_percentage": 94.7,
								"skin_temp_celsius": 34.1
						}
					}
			],
			"next_token": null
	}`)
	})

	ctx := context.Background()
	resp, err := client.Recovery.ListAll(ctx, nil)

	if err != nil {
		t.Fatalf("Recovery.ListAll(): expected nil error, got %#v", err)
	}
	if len(resp.Records) != 2 {
		t.Errorf("Recovery.ListAll(): expected 2 records, got %v", len(resp.Records))
	}
	if resp.NextToken != nil {
		t.Errorf("Recovery.ListAll(): expected next_token nil, got %v", resp.NextToken)
	}
	if resp.Records[0].CycleID != 1 {
		t.Errorf("Recovery.ListAll(): expected record[0] to have CycleID 1, got %v", resp.Records[0].CycleID)
	}
	if resp.Records[1].CycleID != 2 {
		t.Errorf("Recovery.ListAll(): expected record[0] to have CycleID 1, got %v", resp.Records[0].CycleID)
	}
	if resp.Records[0].SleepID != 1 {
		t.Errorf("Recovery.ListAll(): expected record[0] to have SleepID 1, got %v", resp.Records[0].SleepID)
	}
	if resp.Records[1].SleepID != 2 {
		t.Errorf("Recovery.ListAll(): expected record[0] to have SleepID 2, got %v", resp.Records[0].SleepID)
	}
	if resp.Records[0].UserID != 1 {
		t.Errorf("Recovery.ListAll(): expected record[0] to have UserID 1, got %v", resp.Records[0].UserID)
	}
	if resp.Records[1].UserID != 1 {
		t.Errorf("Recovery.ListAll(): expected record[0] to have UserID 1, got %v", resp.Records[0].UserID)
	}
}

func TestRecoveryService_ListAll_with_params(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	date := now()
	mux.HandleFunc("/"+apiVersion+recoveryEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testHttpMethod(t, r, http.MethodGet)
		limit := r.URL.Query().Get("limit")
		nextToken := r.URL.Query().Get("nextToken")
		start := r.URL.Query().Get("start")
		if limit != "1" {
			t.Errorf("Recovery.ListAll(): expected limit == 1, got %v", limit)
		}
		if nextToken != "test_token" {
			t.Errorf("Recovery.ListAll(): expected nextToken == test_token, got %v", nextToken)
		}
		if start != date.Format(time.RFC3339) {
			t.Errorf("Recovery.ListAll(): expected start == %v, got %v", date.Format(time.RFC3339), start)
		}
		fmt.Fprint(w, `
		{
			"records": [
				{
					"cycle_id": 1,
					"sleep_id": 1,
					"user_id": 1,
					"created_at": "2022-11-28T13:29:05.961Z",
					"updated_at": "2022-11-28T13:29:12.588Z",
					"score_state": "SCORED",
					"score": {
							"user_calibrating": false,
							"recovery_score": 69.0,
							"resting_heart_rate": 46.0,
							"hrv_rmssd_milli": 89.81702,
							"spo2_percentage": 93.111115,
							"skin_temp_celsius": 33.4
					}
				}
			],
			"next_token": "test_token"
	}`)
	})

	ctx := context.Background()
	params := RequestParams{Limit: 1, NextToken: "test_token", Start: date}
	resp, err := client.Recovery.ListAll(ctx, &params)

	if err != nil {
		t.Fatalf("Recovery.ListAll(): expected nil error, got %#v", err)
	}
	if len(resp.Records) != 1 {
		t.Errorf("Recovery.ListAll(): expected 1 records, got %v", len(resp.Records))
	}
	if resp.NextToken == nil {
		t.Errorf("Recovery.ListAll(): expected next_token == test_token, got %v", resp.NextToken)
	}
	if resp.Records[0].CycleID != 1 {
		t.Errorf("Recovery.ListAll(): expected record[0] to have CycleID 1, got %v", resp.Records[0].CycleID)
	}
	if resp.Records[0].SleepID != 1 {
		t.Errorf("Recovery.ListAll(): expected record[0] to have SleepID 1, got %v", resp.Records[0].SleepID)
	}
	if resp.Records[0].UserID != 1 {
		t.Errorf("Recovery.ListAll(): expected record[0] to have UserID 1, got %v", resp.Records[0].UserID)
	}
}

func TestRecoveryService_GetOneByCycleId(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/"+apiVersion+cycleEndpoint+"/", func(w http.ResponseWriter, r *http.Request) {
		testHttpMethod(t, r, http.MethodGet)
		urlPath := strings.Split(r.URL.Path, "/")
		if len(urlPath) != 5 {
			t.Errorf("Recovery.GetOneByCycleId(): expected url path of length 4, got %v", len(urlPath))
		}
		if urlPath[len(urlPath)-2] != "1" {
			t.Errorf("Recovery.GetOneByCycleId(): expected id 1, got %v", urlPath[len(urlPath)-2])
		}
		fmt.Fprint(w, `
			{
				"cycle_id": 1,
				"sleep_id": 1,
				"user_id": 1,
				"created_at": "2022-11-28T13:29:05.961Z",
				"updated_at": "2022-11-28T13:29:12.588Z",
				"score_state": "SCORED",
				"score": {
						"user_calibrating": false,
						"recovery_score": 69.0,
						"resting_heart_rate": 46.0,
						"hrv_rmssd_milli": 89.81702,
						"spo2_percentage": 93.111115,
						"skin_temp_celsius": 33.4
				}
			}
	`)
	})

	ctx := context.Background()
	resp, err := client.Recovery.GetOneByCycleId(ctx, 1)

	if err != nil {
		t.Fatalf("Recovery.GetOneByCycleId(): expected nil error, got %#v", err)
	}
	if resp.CycleID != 1 {
		t.Errorf("Recovery.GetOneByCycleId(): expected CycleID 1, got %v", resp.CycleID)
	}
	if resp.SleepID != 1 {
		t.Errorf("Recovery.GetOneByCycleId(): expected SleepID 1, got %v", resp.SleepID)
	}
	if resp.UserID != 1 {
		t.Errorf("Recovery.GetOneByCycleId(): expected UserID 1, got %v", resp.UserID)
	}
}
