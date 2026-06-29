package dbo4spaceus

import (
	"reflect"
	"strings"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

func NewSpaceModuleItemKey[K comparable](spaceID coretypes.SpaceID, moduleID coretypes.ExtID, collection string, itemID K) *dal.Key {
	spaceModuleKey := NewSpaceModuleKey(spaceID, moduleID)
	return dal.NewKeyWithParentAndID(spaceModuleKey, collection, itemID)
}

func NewSpaceModuleItemIncompleteKey[K comparable](spaceID coretypes.SpaceID, moduleID coretypes.ExtID, collection string) *dal.Key {
	spaceModuleKey := NewSpaceModuleKey(spaceID, moduleID)
	var zero K
	idKind := reflect.TypeOf(zero).Kind()
	return dal.NewIncompleteKey(collection, idKind, spaceModuleKey)
}

// specscore: decisions/0002-reserved-extension-space-ids
// Ref→path resolver: the "@{spaceID}" suffix is the sole discriminator between a
// space-bound record and a spaceless system-namespace record (/ext/{ext-id}/...).
// See sneat-specs Decision 0002:
// https://github.com/sneat-co/sneat-specs/blob/main/spec/decisions/0002-reserved-extension-space-ids.md
func NewSpaceModuleItemKeyFromItemRef(spaceID coretypes.SpaceID, itemRef dbo4linkage.ItemRef) *dal.Key {
	// The related ItemID may carry an explicit "@{spaceID}" suffix. The presence
	// or absence of that suffix is the sole discriminator (sneat-specs Decision
	// 0002): a suffix selects the space-bound record in that space; its absence
	// keeps the supplied spaceID (the owner space, or empty for the spaceless
	// system namespace at /ext/{ext-id}/...).
	itemID := itemRef.ItemID
	if i := strings.Index(itemID, dbo4linkage.SpaceItemIDSeparator); i >= 0 {
		spaceID = coretypes.SpaceID(itemID[i+1:])
		itemID = itemID[:i]
	}
	return NewSpaceModuleItemKey(spaceID, itemRef.ExtID, itemRef.Collection, itemID)
}
