package dto4spaceus

import (
	"github.com/strongo/validation"
)

type SpaceItemRequest struct {
	SpaceRequest
	ID string `json:"id"`
}

// Validate returns error if not valid
func (v SpaceItemRequest) Validate() error {
	if v.ID == "" {
		return validation.NewErrRequestIsMissingRequiredField("id")
	}
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	return nil
}
