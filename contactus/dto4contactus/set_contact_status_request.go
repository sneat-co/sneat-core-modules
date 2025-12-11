package dto4contactus

import (
	"strings"

	"github.com/strongo/validation"
)

type SetContactsStatusRequest struct {
	ContactsRequest
	Status string `json:"status"`
}

func (v SetContactsStatusRequest) Validate() error {
	if err := v.ContactsRequest.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(v.Status) == "" {
		return validation.NewErrRequestIsMissingRequiredField("status")
	}
	return nil
}
