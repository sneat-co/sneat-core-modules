package dbo4linkage

import (
	"github.com/sneat-co/sneat-go-core/coretypes"
)

const SpaceItemIDSeparator = "@"

// specscore: decisions/0002-reserved-extension-space-ids
// Ref serialization: omit the "@{spaceID}" suffix for the spaceless system
// namespace. See sneat-specs Decision 0002:
// https://github.com/sneat-co/sneat-specs/blob/main/spec/decisions/0002-reserved-extension-space-ids.md
func NewFullItemRef(extID coretypes.ExtID, collection string, spaceID coretypes.SpaceID, itemID string) ItemRef {
	if itemID == "" {
		panic("itemID is required for a full item reference")
	}
	if spaceID == "" {
		// Spaceless system namespace: no "@{spaceID}" suffix is appended.
		// The record resolves under /ext/{ext-id}/{collection}/{item-id}.
		return newItemRef(extID, collection, itemID)
	}
	return newItemRef(extID, collection, itemID+SpaceItemIDSeparator+string(spaceID))
}
