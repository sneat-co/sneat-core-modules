package dbo4contactus

import (
	"fmt"
	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-core-modules/dbo4all"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/strongo/strongoapp/with"
)

// SpaceContactsCollection defines collection name for space contacts.
// We have `SpaceID` prefix as it can belong only to a single space, and SpaceID is also in record key as prefix.
const SpaceContactsCollection = "contacts"

// ContactDbo belongs only to a single space
type ContactDbo struct {
	//dbmodels.WithSpaceID -- not needed as it's in record key
	//dbmodels.WithUserIDs

	briefs4contactus.ContactBase
	dbo4linkage.WithRelatedAndIDs
	with.CreatedFields
	with.TagsField
	briefs4contactus.WithMultiSpaceContacts[*briefs4contactus.ContactBrief]
	WithInvitesToContactBriefs // dbo4invitus.WithInvites // Invites to become a space member or to connect as a contact
	dbo4all.WithEmails
	dbo4all.WithPhones
}

// Validate returns error if not valid
func (v ContactDbo) Validate() error {
	if err := v.ContactBase.Validate(); err != nil {
		return fmt.Errorf("ContactRecordBase is not valid: %w", err)
	}
	if err := v.CreatedFields.Validate(); err != nil {
		return err
	}
	if err := v.RolesField.Validate(); err != nil {
		return err
	}
	if err := v.TagsField.Validate(); err != nil {
		return err
	}
	if err := v.WithInvitesToContactBriefs.Validate(); err != nil {
		return err
	}
	if err := v.WithRelatedAndIDs.Validate(); err != nil {
		return err
	}
	if err := v.WithPhones.Validate(); err != nil {
		return err
	}
	if err := v.WithEmails.Validate(); err != nil {
		return err
	}
	return nil
}
