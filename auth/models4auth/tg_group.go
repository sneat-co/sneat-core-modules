package models4auth

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
)

const TgGroupKind = "TgGroup"

type TgGroup = record.DataWithID[int64, *TgGroupData]

func NewTgGroup(id int64, data *TgGroupData) TgGroup {
	key := dal.NewKeyWithID(TgGroupKind, id)
	return record.NewDataWithID(id, key, data)
}

//var _ db.EntityHolder = (*TgGroup)(nil)

type TgGroupData struct {
	UserGroupID string `firestore:"userGroupID,omitempty"`
}

//func (TgGroup) Kind() string {
//	return TgGroupKind
//}
//
//func (tgGroup TgGroup) Entity() any {
//	return tgGroup.TgGroupData
//}
//
//func (tgGroup TgGroup) NewEntity() any {
//	return new(TgGroupData)
//}
//
//func (tgGroup *TgGroup) SetEntity(entity any) {
//	if entity == nil {
//		tgGroup.TgGroupData = nil
//	} else {
//		tgGroup.TgGroupData = entity.(*TgGroupData)
//	}
//}
