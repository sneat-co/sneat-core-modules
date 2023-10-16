package dto4schedulus

import (
	"fmt"
	"github.com/sneat-co/sneat-core-modules/schedulus/models4schedulus"
	"github.com/sneat-co/sneat-core-modules/teamus/dto4teamus"
	"github.com/strongo/validation"
)

// CreateHappeningRequest DTO
type CreateHappeningRequest struct {
	dto4teamus.TeamRequest
	Happening *models4schedulus.HappeningBrief `json:"happening"`
}

// Validate returns error if not valid
func (v CreateHappeningRequest) Validate() error {
	if err := v.TeamRequest.Validate(); err != nil {
		return fmt.Errorf("team request is not valid: %w", err)
	}
	if v.Happening == nil {
		return validation.NewErrRequestIsMissingRequiredField("happening")
	}
	if err := v.Happening.Validate(); err != nil {
		return validation.NewErrBadRequestFieldValue("dto", err.Error())
	}
	return nil
}

// CreateHappeningResponse DTO
type CreateHappeningResponse struct {
	ID  string                        `json:"id"`
	Dto models4schedulus.HappeningDto `json:"dto"`
}

// Validate returns error if not valid
func (v CreateHappeningResponse) Validate() error {
	if v.ID == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	return nil
}
