package dbo4linkage

import (
	"github.com/sneat-co/sneat-go-core/coretypes"
)

const SpaceItemIDSeparator = "@"

func NewFullItemRef(extID coretypes.ExtID, collection string, spaceID coretypes.SpaceID, itemID string) ItemRef {
	if spaceID == "" {
		panic("spaceID is required for a full item reference")
	}
	if itemID == "" {
		panic("itemID is required for a full item reference")
	}
	return newItemRef(extID, collection, itemID+SpaceItemIDSeparator+string(spaceID))
}
