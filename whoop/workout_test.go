package whoop

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestWorkoutService_ListAll(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/"+apiVersion+workoutEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testHttpMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `
		{
			"records": [
				{
						"id": 1,
						"user_id": 1,
						"created_at": "2022-12-01T21:26:05.038Z",
						"updated_at": "2022-12-01T21:33:13.930Z",
						"start": "2022-12-01T21:09:54.896Z",
						"end": "2022-12-01T21:26:05.915Z",
						"timezone_offset": "-08:00",
						"sport_id": 27,
						"score_state": "SCORED",
						"score": {
								"strain": 5.1367,
								"average_heart_rate": 110,
								"max_heart_rate": 136,
								"kilojoule": 419.1227,
								"percent_recorded": 100.0,
								"distance_meter": 0.0,
								"altitude_gain_meter": 0.0,
								"altitude_change_meter": 0.0,
								"zone_duration": {
										"zone_zero_milli": 0,
										"zone_one_milli": 297044,
										"zone_two_milli": 650843,
										"zone_three_milli": 24032,
										"zone_four_milli": 0,
										"zone_five_milli": 0
								}
						}
				},
				{
						"id": 2,
						"user_id": 1,
						"created_at": "2022-12-01T21:07:54.985Z",
						"updated_at": "2022-12-01T21:16:16.753Z",
						"start": "2022-12-01T20:11:18.613Z",
						"end": "2022-12-01T21:07:55.673Z",
						"timezone_offset": "-08:00",
						"sport_id": 128,
						"score_state": "SCORED",
						"score": {
								"strain": 13.9132,
								"average_heart_rate": 141,
								"max_heart_rate": 180,
								"kilojoule": 2871.8916,
								"percent_recorded": 100.0,
								"distance_meter": 0.0,
								"altitude_gain_meter": 0.0,
								"altitude_change_meter": 0.0,
								"zone_duration": {
										"zone_zero_milli": 0,
										"zone_one_milli": 55795,
										"zone_two_milli": 160539,
										"zone_three_milli": 2702218,
										"zone_four_milli": 368182,
										"zone_five_milli": 109588
								}
						}
				}
			],
			"next_token": null
	}`)
	})

	ctx := context.Background()
	resp, err := client.Workout.ListAll(ctx, nil)

	if err != nil {
		t.Fatalf("Workout.ListAll(): expected nil error, got %#v", err)
	}
	if len(resp.Records) != 2 {
		t.Errorf("Workout.ListAll(): expected 2 records, got %v", len(resp.Records))
	}
	if resp.NextToken != nil {
		t.Errorf("Workout.ListAll(): expected next_token nil, got %v", resp.NextToken)
	}
	if resp.Records[0].ID != 1 {
		t.Errorf("Workout.ListAll(): expected record[0] to have ID 1, got %v", resp.Records[0].ID)
	}
	if resp.Records[1].ID != 2 {
		t.Errorf("Workout.ListAll(): expected record[1] to have ID 2, got %v", resp.Records[1].ID)
	}
	if resp.Records[0].SportID != 27 {
		t.Errorf("Workout.ListAll(): expected record[0] to have SportID 27, got %v", resp.Records[0].SportID)
	}
	if resp.Records[1].SportID != 128 {
		t.Errorf("Workout.ListAll(): expected record[1] to have SportID 65, got %v", resp.Records[1].SportID)
	}
	if resp.Records[0].SportName != nil && *resp.Records[0].SportName != "Rugby" {
		t.Errorf("Workout.ListAll(): expected record[0] to have SportName Rugby, got %v", resp.Records[0].SportName)
	}
	if resp.Records[1].SportName != nil && *resp.Records[1].SportName != "Stretching" {
		t.Errorf("Workout.ListAll(): expected record[0] to have SportName Stretching, got %v", resp.Records[1].SportName)
	}
	if resp.Records[0].UserID != 1 {
		t.Errorf("Workout.ListAll(): expected record[0] to have UserID 1, got %v", resp.Records[0].UserID)
	}
	if resp.Records[1].UserID != 1 {
		t.Errorf("Workout.ListAll(): expected record[0] to have UserID 1, got %v", resp.Records[0].UserID)
	}
}

func TestWorkoutService_ListAll_with_params(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	date := now()
	mux.HandleFunc("/"+apiVersion+workoutEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testHttpMethod(t, r, http.MethodGet)
		limit := r.URL.Query().Get("limit")
		nextToken := r.URL.Query().Get("nextToken")
		start := r.URL.Query().Get("start")
		if limit != "1" {
			t.Errorf("Workout.ListAll(): expected limit == 1, got %v", limit)
		}
		if nextToken != "test_token" {
			t.Errorf("Workout.ListAll(): expected nextToken == test_token, got %v", nextToken)
		}
		if start != date.Format(time.RFC3339) {
			t.Errorf("Workout.ListAll(): expected start == %v, got %v", date.Format(time.RFC3339), start)
		}
		fmt.Fprint(w, `
		{
			"records": [
				{
						"id": 1,
						"user_id": 1,
						"created_at": "2022-12-01T21:26:05.038Z",
						"updated_at": "2022-12-01T21:33:13.930Z",
						"start": "2022-12-01T21:09:54.896Z",
						"end": "2022-12-01T21:26:05.915Z",
						"timezone_offset": "-08:00",
						"sport_id": 27,
						"score_state": "SCORED",
						"score": {
								"strain": 5.1367,
								"average_heart_rate": 110,
								"max_heart_rate": 136,
								"kilojoule": 419.1227,
								"percent_recorded": 100.0,
								"distance_meter": 0.0,
								"altitude_gain_meter": 0.0,
								"altitude_change_meter": 0.0,
								"zone_duration": {
										"zone_zero_milli": 0,
										"zone_one_milli": 297044,
										"zone_two_milli": 650843,
										"zone_three_milli": 24032,
										"zone_four_milli": 0,
										"zone_five_milli": 0
								}
						}
				}
			],
			"next_token": "test_token"
	}`)
	})

	ctx := context.Background()
	params := RequestParams{Limit: 1, NextToken: "test_token", Start: date}
	resp, err := client.Workout.ListAll(ctx, &params)

	if err != nil {
		t.Fatalf("Workout.ListAll(): expected nil error, got %#v", err)
	}
	if len(resp.Records) != 1 {
		t.Errorf("Workout.ListAll(): expected 1 records, got %v", len(resp.Records))
	}
	if resp.NextToken == nil {
		t.Errorf("Workout.ListAll(): expected next_token == test_token, got %v", resp.NextToken)
	}
	if resp.Records[0].ID != 1 {
		t.Errorf("Workout.ListAll(): expected record[0] to have ID 1, got %v", resp.Records[0].ID)
	}
	if resp.Records[0].SportID != 27 {
		t.Errorf("Workout.ListAll(): expected record[0] to have SportID 27, got %v", resp.Records[0].SportID)
	}
	if resp.Records[0].SportName != nil && *resp.Records[0].SportName != "Rugby" {
		t.Errorf("Workout.ListAll(): expected record[0] to have SportName Rugby, got %v", resp.Records[0].SportName)
	}
	if resp.Records[0].UserID != 1 {
		t.Errorf("Workout.ListAll(): expected record[0] to have UserID 1, got %v", resp.Records[0].UserID)
	}
}

func TestWorkoutService_GetOne(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/"+apiVersion+workoutEndpoint+"/", func(w http.ResponseWriter, r *http.Request) {
		testHttpMethod(t, r, http.MethodGet)
		urlPath := strings.Split(r.URL.Path, "/")
		if len(urlPath) != 5 {
			t.Errorf("Workout.GetOne(): expected url path of length 4, got %v", len(urlPath))
		}
		if urlPath[len(urlPath)-1] != "1" {
			t.Errorf("Workout.GetOne(): expected id 1, got %v", urlPath[len(urlPath)-1])
		}
		if urlPath[len(urlPath)-2] != "workout" {
			t.Errorf("Workout.GetOne(): expected path workout, got %v", urlPath[len(urlPath)-2])
		}
		if urlPath[len(urlPath)-3] != "activity" {
			t.Errorf("Workout.GetOne(): expected path activity, got %v", urlPath[len(urlPath)-3])
		}
		fmt.Fprint(w, `
			{
				"id": 1,
				"user_id": 1,
				"created_at": "2022-12-01T21:26:05.038Z",
				"updated_at": "2022-12-01T21:33:13.930Z",
				"start": "2022-12-01T21:09:54.896Z",
				"end": "2022-12-01T21:26:05.915Z",
				"timezone_offset": "-08:00",
				"sport_id": 27,
				"score_state": "SCORED",
				"score": {
						"strain": 5.1367,
						"average_heart_rate": 110,
						"max_heart_rate": 136,
						"kilojoule": 419.1227,
						"percent_recorded": 100.0,
						"distance_meter": 0.0,
						"altitude_gain_meter": 0.0,
						"altitude_change_meter": 0.0,
						"zone_duration": {
								"zone_zero_milli": 0,
								"zone_one_milli": 297044,
								"zone_two_milli": 650843,
								"zone_three_milli": 24032,
								"zone_four_milli": 0,
								"zone_five_milli": 0
						}
				}
		}
	`)
	})

	ctx := context.Background()
	resp, err := client.Workout.GetOne(ctx, 1)

	if err != nil {
		t.Fatalf("Workout.GetOne(): expected nil error, got %#v", err)
	}
	if resp.ID != 1 {
		t.Errorf("Workout.GetOne(): expected ID 1, got %v", resp.ID)
	}
	if resp.SportID != 27 {
		t.Errorf("Workout.GetOne(): expected SportID 27, got %v", resp.SportID)
	}
	if resp.SportName != nil && *resp.SportName != "Rugby" {
		t.Errorf("Workout.GetOne(): expected SportName Rugby, got %v", resp.SportName)
	}
	if resp.UserID != 1 {
		t.Errorf("Workout.GetOne(): expected UserID 1, got %v", resp.UserID)
	}
}
