package dto4schedulus

import (
	"github.com/strongo/validation"
	"strings"
)

type HappeningContactRequest struct {
	HappeningRequest
	ContactID string `json:"contactID"`
}

func (v HappeningContactRequest) Validate() error {
	if err := v.HappeningRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.ContactID) == "" {
		return validation.NewErrRecordIsMissingRequiredField("contactID")
	}
	return nil
}
