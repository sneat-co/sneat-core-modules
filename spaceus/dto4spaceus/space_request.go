package dto4spaceus

import (
	"strings"

	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/strongo/validation"
)

// NewSpaceRequest creates new space request
func NewSpaceRequest(spaceID coretypes.SpaceID) SpaceRequest {
	if spaceID == "" {
		panic("spaceID is required")
	}
	return SpaceRequest{SpaceID: spaceID}
}

// SpaceRequest request
type SpaceRequest struct {
	SpaceID coretypes.SpaceID `json:"spaceID"`
}

// Validate validates request
func (v SpaceRequest) Validate() error {
	if strings.TrimSpace(string(v.SpaceID)) == "" {
		return validation.NewErrRecordIsMissingRequiredField("spaceID")
	}
	return nil
}
