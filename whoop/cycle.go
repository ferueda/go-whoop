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
type Cycle struct {
	// Unique identifier for the physiological cycle.
	ID int `json:"id"`
	// The WHOOP User for the physiological cycle.
	UserID int `json:"user_id"`
	// The time the cycle was recorded in WHOOP.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// The time the cycle was last updated in WHOOP.
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	// Start time bound of the cycle.
	Start *time.Time `json:"start,omitempty"`
	// End time bound of the cycle. If not present, the user is currently in this cycle.
	End *time.Time `json:"end,omitempty"`
	// The user's timezone offset at the time the cycle was recorded.
	// Follows format for Time Zone Designator (TZD) - '+hh:mm', '-hh:mm', or 'Z'.
	TimezoneOffset *string `json:"timezone_offset,omitempty"`
	// Enum: "SCORED", "PENDING_SCORE", or "UNSCORABLE".
	// SCORED means the cycle was scored and the measurement values will be present.
	// PENDING_SCORE means WHOOP is currently evaluating the cycle.
	// UNSCORABLE means this activity could not be scored for some reason.
	ScoreState *string `json:"score_state,omitempty"`
	// WHOOP's measurements and evaluation of the cycle. Only present if the Cycle State is SCORED
	Score struct {
		Strain           float64 `json:"strain,omitempty"`
		Kilojoule        float64 `json:"kilojoule,omitempty"`
		AverageHeartRate int     `json:"average_heart_rate,omitempty"`
		MaxHeartRate     int     `json:"max_heart_rate,omitempty"`
	} `json:"score,omitempty"`
}

// Get a single physiological cycles by id.
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

type ListAllResponse struct {
	Records   []Cycle `json:"records"`
	NextToken *string `json:"next_token"`
}

// ListAll lists all physiological cycles for the authenticated user.
// Results are paginated and sorted by start time in descending order.
//
// WHOOP API docs: https://developer.whoop.com/api#tag/Cycle/operation/getCycleCollection
func (s *CycleService) ListAll(ctx context.Context, params *RequestParams) (*ListAllResponse, error) {
	u, err := addParams(cycleEndpoint, params)
	if err != nil {
		return nil, err
	}

	var resp ListAllResponse
	if err = s.client.get(ctx, u, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
