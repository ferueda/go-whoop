package whoop

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestSleepService_ListAll(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/"+apiVersion+sleepEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testHttpMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `
		{
			"records": [
				{
						"id": 1,
						"user_id": 1,
						"created_at": "2022-11-28T13:29:05.961Z",
						"updated_at": "2022-11-28T13:29:12.588Z",
						"start": "2022-11-28T08:04:51.371Z",
						"end": "2022-11-28T13:13:55.139Z",
						"timezone_offset": "-08:00",
						"nap": false,
						"score_state": "SCORED",
						"score": {
								"stage_summary": {
										"total_in_bed_time_milli": 18542807,
										"total_awake_time_milli": 1678969,
										"total_no_data_time_milli": 0,
										"total_light_sleep_time_milli": 9644306,
										"total_slow_wave_sleep_time_milli": 3003326,
										"total_rem_sleep_time_milli": 4216206,
										"sleep_cycle_count": 2,
										"disturbance_count": 11
								},
								"sleep_needed": {
										"baseline_milli": 27384483,
										"need_from_sleep_debt_milli": 1936474,
										"need_from_recent_strain_milli": 165237,
										"need_from_recent_nap_milli": 0
								},
								"respiratory_rate": 14.619141,
								"sleep_performance_percentage": 57.0,
								"sleep_consistency_percentage": 72.0,
								"sleep_efficiency_percentage": 90.94544
						}
				},
				{
					"id": 2,
					"user_id": 1,
					"created_at": "2022-11-27T16:34:36.226Z",
					"updated_at": "2022-11-27T16:47:13.425Z",
					"start": "2022-11-27T08:17:45.687Z",
					"end": "2022-11-27T16:33:06.131Z",
					"timezone_offset": "-08:00",
					"nap": false,
					"score_state": "SCORED",
					"score": {
							"stage_summary": {
									"total_in_bed_time_milli": 29732519,
									"total_awake_time_milli": 2847947,
									"total_no_data_time_milli": 0,
									"total_light_sleep_time_milli": 14785608,
									"total_slow_wave_sleep_time_milli": 5198067,
									"total_rem_sleep_time_milli": 6900897,
									"sleep_cycle_count": 5,
									"disturbance_count": 18
							},
							"sleep_needed": {
									"baseline_milli": 27384650,
									"need_from_sleep_debt_milli": 2464618,
									"need_from_recent_strain_milli": 260993,
									"need_from_recent_nap_milli": 0
							},
							"respiratory_rate": 13.911133,
							"sleep_performance_percentage": 89.0,
							"sleep_consistency_percentage": null,
							"sleep_efficiency_percentage": 90.68238
					}
				}
			],
			"next_token": null
	}`)
	})

	ctx := context.Background()
	resp, err := client.Sleep.ListAll(ctx, nil)

	if err != nil {
		t.Fatalf("Sleep.ListAll(): expected nil error, got %#v", err)
	}
	if len(resp.Records) != 2 {
		t.Errorf("Sleep.ListAll(): expected 2 records, got %v", len(resp.Records))
	}
	if resp.NextToken != nil {
		t.Errorf("Sleep.ListAll(): expected next_token nil, got %v", resp.NextToken)
	}
	if resp.Records[0].ID != 1 {
		t.Errorf("Sleep.ListAll(): expected record[0] to have ID 1, got %v", resp.Records[0].ID)
	}
	if resp.Records[1].ID != 2 {
		t.Errorf("Sleep.ListAll(): expected record[0] to have ID 1, got %v", resp.Records[0].ID)
	}
	if resp.Records[0].UserID != 1 {
		t.Errorf("Sleep.ListAll(): expected record[0] to have UserID 1, got %v", resp.Records[0].UserID)
	}
	if resp.Records[1].UserID != 1 {
		t.Errorf("Sleep.ListAll(): expected record[0] to have UserID 1, got %v", resp.Records[0].UserID)
	}
}

func TestSleepService_ListAll_with_params(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	date := now()
	mux.HandleFunc("/"+apiVersion+sleepEndpoint, func(w http.ResponseWriter, r *http.Request) {
		testHttpMethod(t, r, http.MethodGet)
		limit := r.URL.Query().Get("limit")
		nextToken := r.URL.Query().Get("nextToken")
		start := r.URL.Query().Get("start")
		if limit != "1" {
			t.Errorf("Sleep.ListAll(): expected limit == 1, got %v", limit)
		}
		if nextToken != "test_token" {
			t.Errorf("Sleep.ListAll(): expected nextToken == test_token, got %v", nextToken)
		}
		if start != date.Format(time.RFC3339) {
			t.Errorf("Sleep.ListAll(): expected start == %v, got %v", date.Format(time.RFC3339), start)
		}
		fmt.Fprint(w, `
		{
			"records": [
				{
					"id": 1,
					"user_id": 1,
					"created_at": "2022-11-28T13:29:05.961Z",
					"updated_at": "2022-11-28T13:29:12.588Z",
					"start": "2022-11-28T08:04:51.371Z",
					"end": "2022-11-28T13:13:55.139Z",
					"timezone_offset": "-08:00",
					"nap": false,
					"score_state": "SCORED",
					"score": {
							"stage_summary": {
									"total_in_bed_time_milli": 18542807,
									"total_awake_time_milli": 1678969,
									"total_no_data_time_milli": 0,
									"total_light_sleep_time_milli": 9644306,
									"total_slow_wave_sleep_time_milli": 3003326,
									"total_rem_sleep_time_milli": 4216206,
									"sleep_cycle_count": 2,
									"disturbance_count": 11
							},
							"sleep_needed": {
									"baseline_milli": 27384483,
									"need_from_sleep_debt_milli": 1936474,
									"need_from_recent_strain_milli": 165237,
									"need_from_recent_nap_milli": 0
							},
							"respiratory_rate": 14.619141,
							"sleep_performance_percentage": 57.0,
							"sleep_consistency_percentage": 72.0,
							"sleep_efficiency_percentage": 90.94544
					}
				}
			],
			"next_token": "test_token"
	}`)
	})

	ctx := context.Background()
	params := RequestParams{Limit: 1, NextToken: "test_token", Start: date}
	resp, err := client.Sleep.ListAll(ctx, &params)

	if err != nil {
		t.Fatalf("Sleep.ListAll(): expected nil error, got %#v", err)
	}
	if len(resp.Records) != 1 {
		t.Errorf("Sleep.ListAll(): expected 1 records, got %v", len(resp.Records))
	}
	if resp.NextToken == nil {
		t.Errorf("Sleep.ListAll(): expected next_token == test_token, got %v", resp.NextToken)
	}
	if resp.Records[0].ID != 1 {
		t.Errorf("Sleep.ListAll(): expected record[0] to have ID 1, got %v", resp.Records[0].ID)
	}
	if resp.Records[0].UserID != 1 {
		t.Errorf("Sleep.ListAll(): expected record[0] to have UserID 1, got %v", resp.Records[0].UserID)
	}
}

func TestSleepService_GetOne(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/"+apiVersion+sleepEndpoint+"/", func(w http.ResponseWriter, r *http.Request) {
		testHttpMethod(t, r, http.MethodGet)
		urlPath := strings.Split(r.URL.Path, "/")
		if len(urlPath) != 5 {
			t.Errorf("Sleep.GetOne(): expected url path of length 4, got %v", len(urlPath))
		}
		if urlPath[len(urlPath)-1] != "1" {
			t.Errorf("Sleep.GetOne(): expected id 1, got %v", urlPath[len(urlPath)-1])
		}
		fmt.Fprint(w, `
			{
				"id": 1,
				"user_id": 1,
				"created_at": "2022-11-28T13:29:05.961Z",
				"updated_at": "2022-11-28T13:29:12.588Z",
				"start": "2022-11-28T08:04:51.371Z",
				"end": "2022-11-28T13:13:55.139Z",
				"timezone_offset": "-08:00",
				"nap": false,
				"score_state": "SCORED",
				"score": {
						"stage_summary": {
								"total_in_bed_time_milli": 18542807,
								"total_awake_time_milli": 1678969,
								"total_no_data_time_milli": 0,
								"total_light_sleep_time_milli": 9644306,
								"total_slow_wave_sleep_time_milli": 3003326,
								"total_rem_sleep_time_milli": 4216206,
								"sleep_cycle_count": 2,
								"disturbance_count": 11
						},
						"sleep_needed": {
								"baseline_milli": 27384483,
								"need_from_sleep_debt_milli": 1936474,
								"need_from_recent_strain_milli": 165237,
								"need_from_recent_nap_milli": 0
						},
						"respiratory_rate": 14.619141,
						"sleep_performance_percentage": 57.0,
						"sleep_consistency_percentage": 72.0,
						"sleep_efficiency_percentage": 90.94544
				}
			}
	`)
	})

	ctx := context.Background()
	resp, err := client.Sleep.GetOne(ctx, 1)

	if err != nil {
		t.Fatalf("Sleep.GetOne(): expected nil error, got %#v", err)
	}
	if resp.ID != 1 {
		t.Errorf("Sleep.GetOne(): expected ID 1, got %v", resp.ID)
	}
	if resp.UserID != 1 {
		t.Errorf("Sleep.GetOne(): expected UserID 1, got %v", resp.UserID)
	}
}
