package dal4teamus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
)

const Collection = "modules"

func NewTeamModuleKey(teamID, moduleID string) *dal.Key {
	teamKey := NewTeamKey(teamID)
	return dal.NewKeyWithParentAndID(teamKey, Collection, moduleID)
}

func NewTeamModuleEntry[D any](teamID, moduleID string, data D) record.DataWithID[string, D] {
	key := NewTeamModuleKey(teamID, moduleID)
	return record.NewDataWithID(moduleID, key, data)
}
