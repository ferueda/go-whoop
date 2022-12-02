package whoop

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestUserService_GetProfile(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/"+apiVersion+userEndpoint+"/", func(w http.ResponseWriter, r *http.Request) {
		testHttpMethod(t, r, http.MethodGet)
		urlPath := strings.Split(r.URL.Path, "/")
		if len(urlPath) != 5 {
			t.Errorf("User.GetProfile(): expected url path of length 4, got %v", len(urlPath))
		}
		if urlPath[len(urlPath)-1] != "basic" {
			t.Errorf("User.GetProfile(): expected url path to be basic, got %v", urlPath[len(urlPath)-1])
		}
		if urlPath[len(urlPath)-2] != "profile" {
			t.Errorf("User.GetProfile(): expected url path to be profile, got %v", urlPath[len(urlPath)-2])
		}
		if urlPath[len(urlPath)-3] != "user" {
			t.Errorf("User.GetProfile(): expected url path to be user, got %v", urlPath[len(urlPath)-3])
		}
		fmt.Fprint(w, `
			{
				"user_id": 1,
				"email": "test@test.com",
				"first_name": "test_first",
				"last_name": "test_last"
			}
	`)
	})

	ctx := context.Background()
	resp, err := client.User.GetProfile(ctx)

	if err != nil {
		t.Fatalf("User.GetProfile(): expected nil error, got %#v", err)
	}
	if resp.ID != 1 {
		t.Errorf("User.GetProfile(): expected ID 1, got %v", resp.ID)
	}
	if *resp.FirstName != "test_first" {
		t.Errorf("User.GetProfile(): expected FirstName test_first, got %v", resp.FirstName)
	}
	if *resp.LastName != "test_last" {
		t.Errorf("User.GetProfile(): expected FirstName test_last, got %v", resp.LastName)
	}
	if *resp.Email != "test@test.com" {
		t.Errorf("User.GetProfile(): expected Email test@test.com, got %v", resp.Email)
	}
}

func TestUserService_GetBodyMeasurement(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/"+apiVersion+userEndpoint+"/", func(w http.ResponseWriter, r *http.Request) {
		testHttpMethod(t, r, http.MethodGet)
		urlPath := strings.Split(r.URL.Path, "/")
		if len(urlPath) != 5 {
			t.Errorf("User.GetBodyMeasurement(): expected url path of length 4, got %v", len(urlPath))
		}
		if urlPath[len(urlPath)-1] != "body" {
			t.Errorf("User.GetBodyMeasurement(): expected url path to be body, got %v", urlPath[len(urlPath)-1])
		}
		if urlPath[len(urlPath)-2] != "measurement" {
			t.Errorf("User.GetBodyMeasurement(): expected url path to be measurement, got %v", urlPath[len(urlPath)-2])
		}
		if urlPath[len(urlPath)-3] != "user" {
			t.Errorf("User.GetBodyMeasurement(): expected url path to be user, got %v", urlPath[len(urlPath)-3])
		}
		fmt.Fprint(w, `
			{
				"height_meter": 1.8288,
				"weight_kilogram": 90.7185,
				"max_heart_rate": 200
			}
	`)
	})

	ctx := context.Background()
	resp, err := client.User.GetBodyMeasurement(ctx)

	if err != nil {
		t.Fatalf("User.GetBodyMeasurement(): expected nil error, got %#v", err)
	}
	if resp.HeightMeter != 1.8288 {
		t.Errorf("User.GetBodyMeasurement(): expected HeightMeter 1.8288, got %v", resp.HeightMeter)
	}
	if resp.WeightKilogram != 90.7185 {
		t.Errorf("User.GetBodyMeasurement(): expected FirstName 90.7185, got %v", resp.WeightKilogram)
	}
	if resp.MaxHeartRate != 200 {
		t.Errorf("User.GetBodyMeasurement(): expected FirstName MaxHeartRate, got %v", resp.MaxHeartRate)
	}
}
