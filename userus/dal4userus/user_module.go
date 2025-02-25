package dal4userus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/core/coremodels"
	"github.com/sneat-co/sneat-go-core/coretypes"

	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
)

const UserModulesCollection = coremodels.ModulesCollection

func NewUserModuleKey(userID string, moduleID coretypes.ModuleID) *dal.Key {
	userKey := dbo4userus.NewUserKey(userID)
	return dal.NewKeyWithParentAndID(userKey, UserModulesCollection, moduleID)
}

func NewUserModuleEntry[D any](userID string, moduleID coretypes.ModuleID, data D) record.DataWithID[coretypes.ModuleID, D] {
	key := NewUserModuleKey(userID, moduleID)
	return record.NewDataWithID(moduleID, key, data)
}
