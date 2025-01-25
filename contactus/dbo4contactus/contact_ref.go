package dbo4contactus

import (
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
)

func NewContactFullRef(spaceID, contactID string) dbo4linkage.SpaceModuleItemRef {
	return dbo4linkage.NewSpaceModuleItemRef(spaceID, const4contactus.ModuleID, const4contactus.ContactsCollection, contactID)
}
