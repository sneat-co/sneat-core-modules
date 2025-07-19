package dal4userus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/core/coremodels"
	"github.com/sneat-co/sneat-go-core/coretypes"

	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
)

const UserExtCollection = coremodels.ExtCollection

func NewUserExtKey(userID string, extID coretypes.ExtID) *dal.Key {
	userKey := dbo4userus.NewUserKey(userID)
	return dal.NewKeyWithParentAndID(userKey, UserExtCollection, extID)
}

func NewUserExtEntry[D any](userID string, extID coretypes.ExtID, data D) record.DataWithID[coretypes.ExtID, D] {
	key := NewUserExtKey(userID, extID)
	return record.NewDataWithID(extID, key, data)
}
