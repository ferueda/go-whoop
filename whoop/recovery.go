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
	CycleID    int        `json:"cycle_id"`
	SleepID    int        `json:"sleep_id"`
	UserID     int        `json:"user_id"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	ScoreState *string    `json:"score_state,omitempty"`
	Score      struct {
		UserCalibrating  bool    `json:"user_calibrating,omitempty"`
		RecoveryScore    float64 `json:"recovery_score,omitempty"`
		RestingHeartRate float64 `json:"resting_heart_rate,omitempty"`
		HrvRmssdMilli    float64 `json:"hrv_rmssd_milli,omitempty"`
		Spo2Percentage   float64 `json:"spo2_percentage,omitempty"`
		SkinTempCelsius  float64 `json:"skin_temp_celsius,omitempty"`
	} `json:"score,omitempty"`
}

// Get a single recovery record for the given cycle.
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
