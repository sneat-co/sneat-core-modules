package dbo4userus

import (
	"github.com/dal-go/dalgo/dal"
)

// Kind is defining collection name for users records
const Kind = "users"

// NewUserKey creates new user doc ref
func NewUserKey[T comparable](userID T) *dal.Key {
	var zero T
	if userID == zero {
		panic("userID is empty value")
	}
	return dal.NewKeyWithID(Kind, userID)
}

// NewUserKeys creates new api4meetingus doc refs
func NewUserKeys[T comparable](userIDs []T) (userKeys []*dal.Key) {
	userKeys = make([]*dal.Key, len(userIDs))
	for i, id := range userIDs {
		userKeys[i] = NewUserKey(id)
	}
	return userKeys
}
