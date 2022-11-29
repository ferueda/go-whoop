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

type Sleep struct {
	ID             int        `json:"id"`
	UserID         int        `json:"user_id"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	Start          *time.Time `json:"start,omitempty"`
	End            *time.Time `json:"end,omitempty"`
	TimezoneOffset *string    `json:"timezone_offset,omitempty"`
	Nap            bool       `json:"nap,omitempty"`
	ScoreState     *string    `json:"score_state,omitempty"`
	Score          struct {
		StageSummary struct {
			TotalInBedTimeMilli         int `json:"total_in_bed_time_milli,omitempty"`
			TotalAwakeTimeMilli         int `json:"total_awake_time_milli,omitempty"`
			TotalNoDataTimeMilli        int `json:"total_no_data_time_milli,omitempty"`
			TotalLightSleepTimeMilli    int `json:"total_light_sleep_time_milli,omitempty"`
			TotalSlowWaveSleepTimeMilli int `json:"total_slow_wave_sleep_time_milli,omitempty"`
			TotalRemSleepTimeMilli      int `json:"total_rem_sleep_time_milli,omitempty"`
			SleepCycleCount             int `json:"sleep_cycle_count,omitempty"`
			DisturbanceCount            int `json:"disturbance_count,omitempty"`
		} `json:"stage_summary,omitempty"`
		SleepNeeded struct {
			BaselineMilli             int `json:"baseline_milli,omitempty"`
			NeedFromSleepDebtMilli    int `json:"need_from_sleep_debt_milli,omitempty"`
			NeedFromRecentStrainMilli int `json:"need_from_recent_strain_milli,omitempty"`
			NeedFromRecentNapMilli    int `json:"need_from_recent_nap_milli,omitempty"`
		} `json:"sleep_needed,omitempty"`
		RespiratoryRate            float64 `json:"respiratory_rate,omitempty"`
		SleepPerformancePercentage float64 `json:"sleep_performance_percentage,omitempty"`
		SleepConsistencyPercentage float64 `json:"sleep_consistency_percentage,omitempty"`
		SleepEfficiencyPercentage  float64 `json:"sleep_efficiency_percentage,omitempty"`
	} `json:"score,omitempty"`
}

// Get a single sleep record for the specified id.
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
