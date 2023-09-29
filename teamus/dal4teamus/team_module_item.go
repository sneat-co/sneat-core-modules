package dal4teamus

import "github.com/dal-go/dalgo/dal"

func NewTeamModuleItemKey[K comparable](teamID, moduleID, collection string, itemID K) *dal.Key {
	teamModuleKey := NewTeamModuleKey(teamID, moduleID)
	return dal.NewKeyWithParentAndID(teamModuleKey, collection, itemID)
}
