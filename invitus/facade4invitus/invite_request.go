package facade4invitus

import (
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/strongo/validation"
)

type InviteRequest struct {
	dto4spaceus.SpaceRequest
	InviteID string `json:"inviteID"`
	Pin      string `json:"pin"`
}

// Validate validates request
func (v *InviteRequest) Validate() error {
	if err := v.SpaceRequest.Validate(); err != nil {
		return err
	}
	if v.InviteID == "" {
		return validation.NewErrRequestIsMissingRequiredField("inviteID")
	}
	if v.Pin == "" {
		return validation.NewErrRequestIsMissingRequiredField("Pin")
	}
	return nil
}
