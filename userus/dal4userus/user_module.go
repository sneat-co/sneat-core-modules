package dal4userus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/core/coremodels"
	"github.com/sneat-co/sneat-go-core/coretypes"

	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
)

const ExtUserCollection = coremodels.ExtCollection

func NewExtUserKey(userID string, extID coretypes.ModuleID) *dal.Key {
	userKey := dbo4userus.NewUserKey(userID)
	return dal.NewKeyWithParentAndID(userKey, ExtUserCollection, extID)
}

func NewUserModuleEntry[D any](userID string, extID coretypes.ModuleID, data D) record.DataWithID[coretypes.ModuleID, D] {
	key := NewExtUserKey(userID, extID)
	return record.NewDataWithID(extID, key, data)
}
