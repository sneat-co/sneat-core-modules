package dal4userus

import (
	"context"
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-core/coretypes"

	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"time"
)

// TxUpdateUser update user record
var TxUpdateUser = func(
	ctx context.Context,
	transaction dal.ReadwriteTransaction,
	timestamp time.Time,
	userKey *dal.Key,
	data []update.Update,
	opts ...dal.Precondition,
) error {
	if transaction == nil {
		panic("transaction == nil")
	}
	if userKey == nil {
		panic("userKey == nil")
	}
	data = append(data,
		update.ByFieldName("timestamp", timestamp),
	)
	return transaction.Update(ctx, userKey, data, opts...)
}

// TxGetUsers load user records
func TxGetUsers(ctx context.Context, tx dal.ReadwriteTransaction, users []dal.Record) (err error) {
	if len(users) == 0 {
		return nil
	}
	return tx.GetMulti(ctx, users)
}

func GetUserSpaceContactID(ctx context.Context, tx dal.ReadSession, userID string, contactusSpaceEntry dal4contactus.ContactusSpaceEntry) (userContactID string, err error) {

	userContactID, _ = contactusSpaceEntry.Data.GetContactBriefByUserID(userID)

	if userContactID != "" {
		return userContactID, nil
	}

	user := dbo4userus.NewUserEntry(userID)

	if err = GetUser(ctx, tx, user); err != nil || !user.Record.Exists() {
		return "", err
	}

	spaceID := contactusSpaceEntry.Key.Parent().ID.(coretypes.SpaceID)

	userSpaceBrief := user.Data.Spaces[spaceID]

	if userSpaceBrief == nil {
		return "", errors.New("user's space brief is not found in user record")
	}

	if userSpaceBrief.UserContactID == "" {
		return "", errors.New("user's space brief has no value in 'userContactID' field")
	}

	return userSpaceBrief.UserContactID, nil
}
