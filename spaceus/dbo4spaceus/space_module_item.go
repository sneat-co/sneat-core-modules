package dbo4spaceus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
)

func NewSpaceModuleItemKey[K comparable](spaceID, moduleID, collection string, itemID K) *dal.Key {
	spaceModuleKey := NewSpaceModuleKey(spaceID, moduleID)
	return dal.NewKeyWithParentAndID(spaceModuleKey, collection, itemID)
}

func NewSpaceModuleItemKeyFromItemRef(itemRef dbo4linkage.SpaceModuleItemRef) *dal.Key {
	return NewSpaceModuleItemKey(itemRef.Space, itemRef.Module, itemRef.Collection, itemRef.ItemID)
}
