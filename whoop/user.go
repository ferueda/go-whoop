package whoop

import (
	"context"
	"fmt"
)

const (
	userEndpoint = "/user"
)

// UserService handles communication with the User related
// endpoints of the API.
type UserService service

// UserProfile represents a member's user profile.
//
// WHOOP API docs: https://developer.whoop.com/docs/developing/user-data/user
type UserProfile struct {
	ID        int     `json:"user_id"`              // The WHOOP User.
	Email     *string `json:"email,omitempty"`      // User's Email.
	FirstName *string `json:"first_name,omitempty"` // User's First Name.
	LastName  *string `json:"last_name,omitempty"`  // User's Last Name
}

// Body measurements about the user, such as their weight and height.
//
// WHOOP API docs: https://developer.whoop.com/docs/developing/user-data/user#body-measurements
type BodyMeasurement struct {
	HeightMeter    float64 `json:"height_meter,omitempty"`    // User's height in meters.
	WeightKilogram float64 `json:"weight_kilogram,omitempty"` // User's weight in kilograms.
	MaxHeartRate   int     `json:"max_heart_rate,omitempty"`  // The max heart rate WHOOP calculated for the user.
}

// GetProfile retrieves the profile for the authenticated user.
//
// WHOOP API docs: https://developer.whoop.com/api#tag/User/operation/getProfileBasic
func (s *UserService) GetProfile(ctx context.Context) (*UserProfile, error) {
	var profile UserProfile
	u := fmt.Sprintf("%v/%v/%v", userEndpoint, "profile", "basic")
	if err := s.client.get(ctx, u, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

// GetBodyMeasurement retrieves the body measurements for the authenticated user.
//
// WHOOP API docs: https://developer.whoop.com/api#tag/User/operation/getBodyMeasurement
func (s *UserService) GetBodyMeasurement(ctx context.Context) (*BodyMeasurement, error) {
	var bodyMeasurement BodyMeasurement
	u := fmt.Sprintf("%v/%v/%v", userEndpoint, "measurement", "body")
	if err := s.client.get(ctx, u, &bodyMeasurement); err != nil {
		return nil, err
	}
	return &bodyMeasurement, nil
}
