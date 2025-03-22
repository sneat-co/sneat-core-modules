package dbo4spaceus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"reflect"
)

func NewSpaceModuleItemKey[K comparable](spaceID coretypes.SpaceID, moduleID coretypes.ModuleID, collection string, itemID K) *dal.Key {
	spaceModuleKey := NewSpaceModuleKey(spaceID, moduleID)
	return dal.NewKeyWithParentAndID(spaceModuleKey, collection, itemID)
}

func NewSpaceModuleItemIncompleteKey[K comparable](spaceID coretypes.SpaceID, moduleID coretypes.ModuleID, collection string) *dal.Key {
	spaceModuleKey := NewSpaceModuleKey(spaceID, moduleID)
	var zero K
	idKind := reflect.TypeOf(zero).Kind()
	return dal.NewIncompleteKey(collection, idKind, spaceModuleKey)
}

func NewSpaceModuleItemKeyFromItemRef(itemRef dbo4linkage.SpaceModuleItemRef) *dal.Key {
	return NewSpaceModuleItemKey(itemRef.Space, itemRef.Module, itemRef.Collection, itemRef.ItemID)
}
