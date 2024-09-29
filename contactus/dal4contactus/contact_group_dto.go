package dal4contactus

import (
	briefs4contactus2 "github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
)

type ContactGroupDto struct {
	briefs4contactus2.ContactGroupBrief
	briefs4contactus2.WithMultiSpaceContacts[*briefs4contactus2.ContactBrief]
}

func (v *ContactGroupDto) Validate() error {
	if err := v.ContactGroupBrief.Validate(); err != nil {
		return err
	}
	if err := v.WithMultiSpaceContacts.Validate(); err != nil {
		return err
	}
	return nil
}
