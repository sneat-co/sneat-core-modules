package dto4spaceus

import (
	"github.com/strongo/validation"
	"strings"
)

// NewSpaceRequest creates new space request
func NewSpaceRequest(spaceID string) SpaceRequest {
	if spaceID == "" {
		panic("spaceID is required")
	}
	return SpaceRequest{SpaceID: spaceID}
}

// SpaceRequest request
type SpaceRequest struct {
	SpaceID string `json:"spaceID"`
}

// Validate validates request
func (v SpaceRequest) Validate() error {
	if strings.TrimSpace(v.SpaceID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("spaceID")
	}
	return nil
}
