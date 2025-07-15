package facade4invitus

import (
	"github.com/strongo/validation"
)

type InviteRequest struct {
	InviteID string `json:"inviteID"`
	Pin      string `json:"pin"`
}

// Validate validates request
func (v *InviteRequest) Validate() error {
	if v.InviteID == "" {
		return validation.NewErrRequestIsMissingRequiredField("inviteID")
	}
	if v.Pin == "" {
		return validation.NewErrRequestIsMissingRequiredField("Pin")
	}
	return nil
}
