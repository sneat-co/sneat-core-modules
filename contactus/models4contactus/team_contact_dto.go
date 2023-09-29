package models4contactus

import (
	"fmt"
	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-core-modules/invitus/models4invitus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
)

// TeamContactsCollection defines  collection name for team contacts.
// We have `Team` prefix as it can belong only to single team
// and TeamID is also in record key as prefix.
const TeamContactsCollection = "contacts"

// ContactDto belongs only to single team
type ContactDto struct {
	//dbmodels.WithTeamID -- not needed as it's in record key
	//dbmodels.WithUserIDs
	briefs4contactus.ContactBase
	dbmodels.WithTags
	briefs4contactus.WithMultiTeamContacts[*briefs4contactus.ContactBrief]
	models4invitus.WithInvites // Invites to become a team member
}

// Validate returns error if not valid
func (v ContactDto) Validate() error {
	if err := v.ContactBase.Validate(); err != nil {
		return fmt.Errorf("ContactRecordBase is not valid: %w", err)
	}
	if err := v.WithRoles.Validate(); err != nil {
		return err
	}
	if err := v.WithTags.Validate(); err != nil {
		return err
	}
	if err := v.WithInvites.Validate(); err != nil {
		return err
	}
	return nil
}
