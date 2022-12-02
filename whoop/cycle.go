package whoop

import (
	"context"
	"fmt"
	"time"
)

const (
	cycleEndpoint = "/cycle"
)

// CycleService handles communication with the Cycle related
// endpoints of the API.
type CycleService service

// Cycle represents a member's activity in the context of a Physiological Cycle.
//
// WHOOP API docs: https://developer.whoop.com/docs/developing/user-data/cycle
type Cycle struct {
	ID     int `json:"id"`      // Unique identifier for the physiological cycle.
	UserID int `json:"user_id"` // The User for the physiological cycle.

	CreatedAt *time.Time `json:"created_at,omitempty"` // Time the cycle was recorded.
	UpdatedAt *time.Time `json:"updated_at,omitempty"` // Time the cycle was last updated.
	Start     *time.Time `json:"start,omitempty"`      // Start time bound of the cycle.
	End       *time.Time `json:"end,omitempty"`        // End time bound of the cycle. If not present, the user is currently in this cycle.

	TimezoneOffset *string `json:"timezone_offset,omitempty"` // Timezone offset at the time the cycle was recorded.
	ScoreState     *string `json:"score_state,omitempty"`     // "SCORED", "PENDING_SCORE", or "UNSCORABLE".

	Score struct {
		Strain           float64 `json:"strain,omitempty"`             // Level of strain for the user. Scored from 0 to 21.
		Kilojoule        float64 `json:"kilojoule,omitempty"`          // Kilojoules the user expended during the cycle.
		AverageHeartRate float64 `json:"average_heart_rate,omitempty"` // The user's average heart rate during the cycle.
		MaxHeartRate     float64 `json:"max_heart_rate,omitempty"`     // The user's max heart rate during the cycle.
	} `json:"score,omitempty"`
}

// GetOne retrieves a single physiological cycle record for the specified id.
//
// WHOOP API docs: https://developer.whoop.com/api#tag/Cycle/operation/getCycleById
func (s *CycleService) GetOne(ctx context.Context, id int) (*Cycle, error) {
	var cycle Cycle
	u := fmt.Sprintf("%v/%v", cycleEndpoint, id)
	if err := s.client.get(ctx, u, &cycle); err != nil {
		return nil, err
	}
	return &cycle, nil
}

type CycleListAllResp struct {
	Records   []Cycle `json:"records"`
	NextToken *string `json:"next_token"`
}

// ListAll lists all physiological cycle records for the authenticated user.
// Results are paginated and sorted by start time in descending order.
//
// WHOOP API docs: https://developer.whoop.com/api#tag/Cycle/operation/getCycleCollection
func (s *CycleService) ListAll(ctx context.Context, params *RequestParams) (*CycleListAllResp, error) {
	u, err := addParams(cycleEndpoint, params)
	if err != nil {
		return nil, err
	}

	var resp CycleListAllResp
	if err = s.client.get(ctx, u, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
