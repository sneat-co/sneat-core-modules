package dbo4contactus

import (
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

func NewContactFullRef(spaceID coretypes.SpaceID, contactID string) dbo4linkage.SpaceModuleItemRef {
	return dbo4linkage.NewSpaceModuleItemRef(spaceID, const4contactus.ModuleID, const4contactus.ContactsCollection, contactID)
}
