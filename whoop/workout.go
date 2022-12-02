package whoop

import (
	"context"
	"fmt"
	"time"
)

const (
	workoutEndpoint = "/activity/workout"
)

// WorkoutService handles communication with the Workout related
// endpoints of the API.
type WorkoutService service

// Workout represents a workout activity for a given user.
//
// WHOOP API docs: https://developer.whoop.com/docs/developing/user-data/workout
type Workout struct {
	ID     int `json:"id"`      // Unique identifier for the workout activity.
	UserID int `json:"user_id"` // The User for the workout activity.

	CreatedAt *time.Time `json:"created_at,omitempty"` // Time the workout was recorded.
	UpdatedAt *time.Time `json:"updated_at,omitempty"` // Time the workout was last updated.
	Start     *time.Time `json:"start,omitempty"`      // Start time bound of the workout.
	End       *time.Time `json:"end,omitempty"`        // End time bound of the workout.

	TimezoneOffset *string `json:"timezone_offset,omitempty"` // // Timezone offset at the time the workout was recorded.
	SportID        int     `json:"sport_id,omitempty"`        // ID of the Sport performed during the workout
	SportName      *string `json:"sport_name,omitempty"`      // Name of the WHOOP Sport performed during the workout
	ScoreState     *string `json:"score_state,omitempty"`     // "SCORED", "PENDING_SCORE", or "UNSCORABLE".

	Score struct {
		Strain              float64 `json:"strain,omitempty"`                // Level of strain of the workout. Scored from 0 to 21.
		AverageHeartRate    int     `json:"average_heart_rate,omitempty"`    // User's average heart rate during the workout.
		MaxHeartRate        int     `json:"max_heart_rate,omitempty"`        // User's max heart rate during the workout.
		Kilojoule           float64 `json:"kilojoule,omitempty"`             // Kilojoules expended during the workout.
		PercentRecorded     float64 `json:"percent_recorded,omitempty"`      // Percentage of heart rate recorded during the workout.
		DistanceMeter       float64 `json:"distance_meter,omitempty"`        // Distance travelled during the workout.
		AltitudeGainMeter   float64 `json:"altitude_gain_meter,omitempty"`   // Altitude gained during the workout.
		AltitudeChangeMeter float64 `json:"altitude_change_meter,omitempty"` // Altitude difference between start and end points of the workout.

		ZoneDuration struct {
			ZoneZeroMilli  int `json:"zone_zero_milli,omitempty"`  // Time spent with Heart Rate lower than Zone One [0-50%).
			ZoneOneMilli   int `json:"zone_one_milli,omitempty"`   // Time spent in Heart Rate Zone One [50-60%)
			ZoneTwoMilli   int `json:"zone_two_milli,omitempty"`   // Time spent in Heart Rate Zone Two [60-70%).
			ZoneThreeMilli int `json:"zone_three_milli,omitempty"` // Time spent in Heart Rate Zone Three [70-80%).
			ZoneFourMilli  int `json:"zone_four_milli,omitempty"`  // Time spent in Heart Rate Zone Four [80-90%).
			ZoneFiveMilli  int `json:"zone_five_milli,omitempty"`  // Time spent in Heart Rate Zone Five [90-100%).
		} `json:"zone_duration,omitempty"`
	} `json:"score,omitempty"`
}

// GetOne retrieves a single workout record for the specified id.
//
// WHOOP API docs: https://developer.whoop.com/api#tag/Workout/operation/getWorkoutById
func (s *WorkoutService) GetOne(ctx context.Context, id int) (*Workout, error) {
	var workout Workout
	u := fmt.Sprintf("%v/%v", workoutEndpoint, id)
	if err := s.client.get(ctx, u, &workout); err != nil {
		return nil, err
	}
	if val, ok := Sports[workout.SportID]; ok {
		workout.SportName = &val
	}
	return &workout, nil
}

type WorkoutListAllResp struct {
	Records   []Workout `json:"records"`
	NextToken *string   `json:"next_token"`
}

// ListAll lists all workout records for the authenticated user.
// Results are paginated and sorted by start time in descending order.
//
// WHOOP API docs: https://developer.whoop.com/api#tag/Workout/operation/getWorkoutCollection
func (s *WorkoutService) ListAll(ctx context.Context, params *RequestParams) (*WorkoutListAllResp, error) {
	u, err := addParams(workoutEndpoint, params)
	if err != nil {
		return nil, err
	}

	var resp WorkoutListAllResp
	if err = s.client.get(ctx, u, &resp); err != nil {
		return nil, err
	}
	for i := range resp.Records {
		if val, ok := Sports[resp.Records[i].SportID]; ok {
			resp.Records[i].SportName = &val
		}
	}
	return &resp, nil
}

// WHOOP sports mapping id to sport name
var Sports map[int]string = map[int]string{
	-1:  "Activity",
	0:   "Running",
	1:   "Cycling",
	16:  "Baseball",
	17:  "Basketball",
	18:  "Rowing",
	19:  "Fencing",
	20:  "Field Hockey",
	21:  "Football",
	22:  "Golf",
	24:  "Ice Hockey",
	25:  "Lacrosse",
	27:  "Rugby",
	28:  "Sailing",
	29:  "Skiing",
	30:  "Soccer",
	31:  "Softball",
	32:  "Squash",
	33:  "Swimming",
	34:  "Tennis",
	35:  "Track & Field",
	36:  "Volleyball",
	37:  "Water Polo",
	38:  "Wrestling",
	39:  "Boxing",
	42:  "Dance",
	43:  "Pilates",
	44:  "Yoga",
	45:  "Weightlifting",
	47:  "Cross Country Skiing",
	48:  "Functional Fitness",
	49:  "Duathlon",
	51:  "Gymnastics",
	52:  "Hiking/Rucking",
	53:  "Horseback Riding",
	55:  "Kayaking",
	56:  "Martial Arts",
	57:  "Mountain Biking",
	59:  "Powerlifting",
	60:  "Rock Climbing",
	61:  "Paddleboarding",
	62:  "Triathlon",
	63:  "Walking",
	64:  "Surfing",
	65:  "Elliptical",
	66:  "Stairmaster",
	70:  "Meditation",
	71:  "Other",
	73:  "Diving",
	74:  "Operations - Tactical",
	75:  "Operations - Medical",
	76:  "Operations - Flying",
	77:  "Operations - Water",
	82:  "Ultimate",
	83:  "Climber",
	84:  "Jumping Rope",
	85:  "Australian Football",
	86:  "Skateboarding",
	87:  "Coaching",
	88:  "Ice Bath",
	89:  "Commuting",
	90:  "Gaming",
	91:  "Snowboarding",
	92:  "Motocross",
	93:  "Caddying",
	94:  "Obstacle Course Racing",
	95:  "Motor Racing",
	96:  "HIIT",
	97:  "Spin",
	98:  "Jiu Jitsu",
	99:  "Manual Labor",
	100: "Cricket",
	101: "Pickleball",
	102: "Inline Skating",
	103: "Box Fitness",
	104: "Spikeball",
	105: "Wheelchair Pushing",
	106: "Paddle Tennis",
	107: "Barre",
	108: "Stage Performance",
	109: "High Stress Work",
	110: "Parkour",
	111: "Gaelic Football",
	112: "Hurling/Camogie",
	113: "Circus Arts",
	121: "Massage Therapy",
	125: "Watching Sports",
	126: "Assault Bike",
	127: "Kickboxing",
	128: "Stretching",
}
