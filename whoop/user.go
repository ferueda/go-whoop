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
	Email     *string `json:"email,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
}

type BodyMeasurement struct {
	HeightMeter    float64 `json:"height_meter,omitempty"`
	WeightKilogram float64 `json:"weight_kilogram,omitempty"`
	MaxHeartRate   int     `json:"max_heart_rate,omitempty"`
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
