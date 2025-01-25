package dbo4spaceus

import (
	"github.com/dal-go/dalgo/record"
)

type SpaceEntry = record.DataWithID[string, *SpaceDbo]

func NewSpaceEntry(id string) (space SpaceEntry) {
	spaceDto := new(SpaceDbo)
	return NewSpaceEntryWithDbo(id, spaceDto)
}

func NewSpaceEntryWithDbo(id string, dbo *SpaceDbo) (space SpaceEntry) {
	if dbo == nil {
		panic("required parameter dbo is nil")
	}
	space = record.NewDataWithID(id, NewSpaceKey(id), dbo)
	space.ID = id
	return
}
