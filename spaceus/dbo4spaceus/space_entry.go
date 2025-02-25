package dbo4spaceus

import (
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

type SpaceEntry = record.DataWithID[coretypes.SpaceID, *SpaceDbo]

func NewSpaceEntry(id coretypes.SpaceID) (space SpaceEntry) {
	spaceDto := new(SpaceDbo)
	return NewSpaceEntryWithDbo(id, spaceDto)
}

func NewSpaceEntryWithDbo(id coretypes.SpaceID, dbo *SpaceDbo) (space SpaceEntry) {
	if dbo == nil {
		panic("required parameter dbo is nil")
	}
	space = record.NewDataWithID(id, NewSpaceKey(id), dbo)
	space.ID = id
	return
}
