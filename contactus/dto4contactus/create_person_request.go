package dto4contactus

import (
	"github.com/sneat-co/sneat-go-core/modules/contactus/briefs4contactus"
	"github.com/strongo/validation"
)

// CreatePersonRequest - base for CreateMemberRequest & facade4contactus.CreateTeamContactRequest
type CreatePersonRequest struct {
	briefs4contactus.ContactBase
	Relationship string `json:"relationship"` // Related to creator
	Message      string `json:"message"`
}

// Validate returns error if not valid
func (v CreatePersonRequest) Validate() error {
	if err := v.ContactBase.Validate(); err != nil {
		return validation.NewBadRequestError(err)
	}
	return nil
}
