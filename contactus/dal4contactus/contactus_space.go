package dal4contactus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-go-core/coretypes"

	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
)

type ContactusSpaceEntry = record.DataWithID[coretypes.ModuleID, *dbo4contactus.ContactusSpaceDbo]

func NewContactusSpaceKey(spaceID coretypes.SpaceID) *dal.Key {
	return dbo4spaceus.NewSpaceModuleKey(spaceID, const4contactus.ModuleID)
}

func NewContactusSpaceEntry(spaceID coretypes.SpaceID) ContactusSpaceEntry {
	return NewContactusSpaceEntryWithData(spaceID, new(dbo4contactus.ContactusSpaceDbo))
}

func NewContactusSpaceEntryWithData(spaceID coretypes.SpaceID, data *dbo4contactus.ContactusSpaceDbo) ContactusSpaceEntry {
	return dbo4spaceus.NewSpaceModuleEntry(spaceID, const4contactus.ModuleID, data)
}
