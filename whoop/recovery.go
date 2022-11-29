package whoop

import (
	"context"
	"fmt"
	"time"
)

const (
	recoveryEndpoint = "/recovery"
)

// RecoveryService handles communication with the Reovery related
// endpoints of the API.
type RecoveryService service

// Recovery represents a member's recovery score in the context of a Physiological Cycle.
//
// WHOOP API docs: https://developer.whoop.com/docs/developing/user-data/recovery
type Recovery struct {
	// Unique identifier for the sleep activity
	CycleID int `json:"cycle_id"`
	// The WHOOP User who performed the sleep activity
	SleepID int `json:"sleep_id"`
	// The WHOOP User for the recovery
	UserID int `json:"user_id"`
	// The time the recovery was recorded in WHOOP
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// The time the recovery was last updated in WHOOP
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	// Enum: "SCORED", "PENDING_SCORE", or "UNSCORABLE".
	// SCORED means the cycle was scored and the measurement values will be present.
	// PENDING_SCORE means WHOOP is currently evaluating the cycle.
	// UNSCORABLE means this activity could not be scored for some reason.
	ScoreState *string `json:"score_state,omitempty"`
	// WHOOP's measurements and evaluation of the recovery. Only present if the Recovery State is SCORED
	Score struct {
		// True if the user is still calibrating and not enough data is available in WHOOP to provide an accurate recovery.
		UserCalibrating bool `json:"user_calibrating,omitempty"`
		// Percentage (0-100%) that reflects how well prepared the user's body is to take on Strain.
		// The Recovery score is a measure of the user body's "return to baseline" after a stressor.
		RecoveryScore float64 `json:"recovery_score,omitempty"`
		// The user's resting heart rate.
		RestingHeartRate float64 `json:"resting_heart_rate,omitempty"`
		// The user's Heart Rate Variability measured using Root Mean Square of Successive Differences (RMSSD), in milliseconds.
		HrvRmssdMilli float64 `json:"hrv_rmssd_milli,omitempty"`
		// The percentage of oxygen in the user's blood. Only present if the user is on 4.0 or greater.
		Spo2Percentage float64 `json:"spo2_percentage,omitempty"`
		// The user's skin temperature, in Celsius. Only present if the user is on 4.0 or greater.
		SkinTempCelsius float64 `json:"skin_temp_celsius,omitempty"`
	} `json:"score,omitempty"`
}

// Get a single recovery record for the specified cycle id.
//
// WHOOP API docs: https://developer.whoop.com/api#tag/Cycle/operation/getCycleById
func (s *RecoveryService) GetOneByCycleId(ctx context.Context, id int) (*Recovery, error) {
	var recovery Recovery
	u := fmt.Sprintf("%v/%v%v", cycleEndpoint, id, recoveryEndpoint)
	if err := s.client.get(ctx, u, &recovery); err != nil {
		return nil, err
	}
	return &recovery, nil
}

type RecoveryListAllResp struct {
	Records   []Recovery `json:"records"`
	NextToken *string    `json:"next_token"`
}

// ListAll lists all recovery records for the authenticated user.
// Results are paginated and sorted by start time in descending order.
//
// WHOOP API docs: https://developer.whoop.com/api#tag/Recovery/operation/getRecoveryCollection
func (s *RecoveryService) ListAll(ctx context.Context, params *RequestParams) (*RecoveryListAllResp, error) {
	u, err := addParams(recoveryEndpoint, params)
	if err != nil {
		return nil, err
	}

	var resp RecoveryListAllResp
	if err = s.client.get(ctx, u, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
