package dbo4contactus

import (
	briefs4contactus2 "github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
)

type ContactusSpaceDbo struct {
	TotalContactsCountByStatus map[string]int `json:"totalContactsCountByStatus,omitempty" firestore:"totalContactsCountByStatus,omitempty"`
	briefs4contactus2.WithSingleSpaceContactsWithoutContactIDs[*briefs4contactus2.ContactBrief]
}

func (v *ContactusSpaceDbo) Validate() error {
	return v.WithSingleSpaceContactsWithoutContactIDs.Validate()
}
