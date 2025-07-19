package dbo4spaceus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/core/coremodels"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

const SpaceModulesCollection = coremodels.ExtCollection

func NewSpaceModuleKey(spaceID coretypes.SpaceID, moduleID coretypes.ExtID) *dal.Key {
	spaceKey := NewSpaceKey(coretypes.SpaceID(spaceID))
	return dal.NewKeyWithParentAndID(spaceKey, SpaceModulesCollection, moduleID)
}

func NewSpaceModuleEntry[D any](spaceID coretypes.SpaceID, moduleID coretypes.ExtID, data D) record.DataWithID[coretypes.ExtID, D] {
	key := NewSpaceModuleKey(spaceID, moduleID)
	return record.NewDataWithID(moduleID, key, data)
}
