package whoop

import (
	"context"
	"fmt"
	"time"
)

const (
	sleepEndpoint = "/activity/sleep"
)

// SleepService handles communication with the Sleep related
// endpoints of the API.
type SleepService service

// Sleep represents a sleep performance for a given user.
//
// WHOOP API docs: https://developer.whoop.com/docs/developing/user-data/sleep
type Sleep struct {
	ID     int `json:"id"`      // Unique identifier for the sleep activity
	UserID int `json:"user_id"` // User who performed the sleep activity

	CreatedAt *time.Time `json:"created_at,omitempty"` // Time the cycle was recorded.
	UpdatedAt *time.Time `json:"updated_at,omitempty"` // Time the cycle was last updated.
	Start     *time.Time `json:"start,omitempty"`      // Start time bound of the cycle.
	End       *time.Time `json:"end,omitempty"`        // End time bound of the cycle. If not present, the user is currently in this cycle.

	TimezoneOffset *string `json:"timezone_offset,omitempty"` // // Timezone offset at the time the cycle was recorded.
	Nap            bool    `json:"nap,omitempty"`             // If true, this sleep activity was a nap for the user.
	ScoreState     *string `json:"score_state,omitempty"`     // "SCORED", "PENDING_SCORE", or "UNSCORABLE".

	Score struct {
		StageSummary struct {
			TotalInBedTimeMilli         int `json:"total_in_bed_time_milli,omitempty"`          // Total time the user spent in bed, in milliseconds.
			TotalAwakeTimeMilli         int `json:"total_awake_time_milli,omitempty"`           // Total time the user spent awake, in milliseconds.
			TotalNoDataTimeMilli        int `json:"total_no_data_time_milli,omitempty"`         // Total time WHOOP did not receive data from the user during the sleep, in milliseconds.
			TotalLightSleepTimeMilli    int `json:"total_light_sleep_time_milli,omitempty"`     // Total time the user spent in light sleep, in milliseconds.
			TotalSlowWaveSleepTimeMilli int `json:"total_slow_wave_sleep_time_milli,omitempty"` // Total time the user spent in Slow Wave Sleep, in milliseconds.
			TotalRemSleepTimeMilli      int `json:"total_rem_sleep_time_milli,omitempty"`       // Total time the user spent in Rapid Eye Movement (REM) sleep, in milliseconds.
			SleepCycleCount             int `json:"sleep_cycle_count,omitempty"`                // Number of sleep cycles during the user's sleep.
			DisturbanceCount            int `json:"disturbance_count,omitempty"`                // Number of times the user was disturbed during sleep
		} `json:"stage_summary,omitempty"`
		SleepNeeded struct {
			BaselineMilli             int `json:"baseline_milli,omitempty"`                // Amount of sleep a user needed based on historical trends.
			NeedFromSleepDebtMilli    int `json:"need_from_sleep_debt_milli,omitempty"`    // Difference between the amount of sleep the user's body required and the amount the user actually got.
			NeedFromRecentStrainMilli int `json:"need_from_recent_strain_milli,omitempty"` // Additional sleep need accrued based on the user's strain.
			NeedFromRecentNapMilli    int `json:"need_from_recent_nap_milli,omitempty"`    // Reduction in sleep need accrued based on the user's recent nap activity (negative value or zero).
		} `json:"sleep_needed,omitempty"`
		RespiratoryRate            float64 `json:"respiratory_rate,omitempty"`             // User's respiratory rate during the sleep.
		SleepPerformancePercentage float64 `json:"sleep_performance_percentage,omitempty"` // Percentage of time user is asleep over the amount of sleep the user needed.
		SleepConsistencyPercentage float64 `json:"sleep_consistency_percentage,omitempty"` // Percentage of how similar this sleep and wake times compared to the previous day.
		SleepEfficiencyPercentage  float64 `json:"sleep_efficiency_percentage,omitempty"`  // Percentage of time user spends in bed that user is actually asleep.
	} `json:"score,omitempty"`
}

// GetOne retrieves a single sleep record for the specified id.
//
// WHOOP API docs: https://developer.whoop.com/api#tag/Sleep/operation/getSleepById
func (s *SleepService) GetOne(ctx context.Context, id int) (*Sleep, error) {
	var sleep Sleep
	u := fmt.Sprintf("%v/%v", sleepEndpoint, id)
	if err := s.client.get(ctx, u, &sleep); err != nil {
		return nil, err
	}
	return &sleep, nil
}

type SleepListAllResp struct {
	Records   []Sleep `json:"records"`
	NextToken *string `json:"next_token"`
}

// ListAll lists all sleep records for the authenticated user.
// Results are paginated and sorted by start time in descending order.
//
// WHOOP API docs: https://developer.whoop.com/api#tag/Sleep/operation/getSleepCollection
func (s *SleepService) ListAll(ctx context.Context, params *RequestParams) (*SleepListAllResp, error) {
	u, err := addParams(sleepEndpoint, params)
	if err != nil {
		return nil, err
	}

	var resp SleepListAllResp
	if err = s.client.get(ctx, u, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
