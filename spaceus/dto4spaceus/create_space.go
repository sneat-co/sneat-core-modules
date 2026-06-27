package dto4spaceus

import (
	"strings"

	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

var _ facade.Request = (*CreateSpaceRequest)(nil)

// CreateSpaceRequest request
type CreateSpaceRequest struct {
	Type  coretypes.SpaceType `json:"type"`
	Title string              `json:"title,omitempty"`
}

// Validate validates request
func (request *CreateSpaceRequest) Validate() error {
	if strings.TrimSpace(string(request.Type)) == "" {
		return validation.NewErrRecordIsMissingRequiredField("type")
	}
	if request.Type == coretypes.SpaceTypeSystem {
		return validation.NewErrBadRequestFieldValue("type", "system spaces are provisioned by the platform only")
	}
	if request.Type != coretypes.SpaceTypeFamily &&
		request.Type != coretypes.SpaceTypePrivate &&
		strings.TrimSpace(request.Title) == "" {
		return validation.NewErrRecordIsMissingRequiredField("title")
	}
	return nil
}

// CreateSpaceResponse response
type CreateSpaceResponse = SpaceResponse
