package models4auth

import (
	"github.com/dal-go/dalgo/record"
	"github.com/strongo/strongoapp/appuser"
)

const (
	UserVkKind = "UserVk"
)

type UserVkEntity struct {
	appuser.OwnedByUserWithID
	FirstName  string
	LastName   string
	ScreenName string
	Nickname   string
	//FriendIDs []int64 `firestore:",omitempty"`
}

type UserVk struct {
	record.WithID[int]
	*UserVkEntity
}

//var _ db.EntityHolder = (*UserVk)(nil)

func (UserVk) Kind() string {
	return UserVkKind
}

func (u UserVk) Entity() any {
	return u.UserVkEntity
}

func (UserVk) NewEntity() any {
	return new(UserVkEntity)
}

func (u *UserVk) SetEntity(entity any) {
	if entity == nil {
		u.UserVkEntity = nil
	} else {
		u.UserVkEntity = entity.(*UserVkEntity)
	}
}
