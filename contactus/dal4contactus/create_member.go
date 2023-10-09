package dal4contactus

import (
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-core-modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/validation"
)

var _ facade.Request = (*CreateMemberRequest)(nil)

// CreateMemberRequest request
type CreateMemberRequest struct {
	dto4teamus.TeamRequest
	dto4contactus.CreatePersonRequest

	RelatedTo *dto4contactus.RelatedToRequest `json:"relatedTo,omitempty"`

	Message string `json:"message"`
}

// Validate validates request
func (v *CreateMemberRequest) Validate() error {
	if err := v.TeamRequest.Validate(); err != nil {
		return err
	}
	if err := v.CreatePersonRequest.Validate(); err != nil {
		return err
	}
	if err := v.RelatedTo.Validate(); err != nil {
		return validation.NewErrBadRequestFieldValue("relatedTo", err.Error())
	}
	// Validate relationship
	return nil
}
