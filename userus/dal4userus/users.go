package dal4userus

import (
	"context"
	"errors"
	"time"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactusmodels/briefs4contactus"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
)

// contactusSpaceContactsReader is the subset of the contactus space module data that
// GetUserSpaceContactID needs, defined here so userus does not depend on contactus DAL types.
type contactusSpaceContactsReader interface {
	GetContactBriefByUserID(userID string) (string, *briefs4contactus.ContactBrief)
}

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

func GetUserSpaceContactID(ctx context.Context, tx dal.ReadSession, userID string, spaceID string, spaceContacts contactusSpaceContactsReader) (userContactID string, err error) {

	userContactID, _ = spaceContacts.GetContactBriefByUserID(userID)

	if userContactID != "" {
		return userContactID, nil
	}

	user := dbo4userus.NewUserEntry(userID)

	if err = GetUser(ctx, tx, user); err != nil || !user.Record.Exists() {
		return "", err
	}

	userSpaceBrief := user.Data.Spaces[spaceID]

	if userSpaceBrief == nil {
		return "", errors.New("user's space brief is not found in user record")
	}

	if userSpaceBrief.UserContactID == "" {
		return "", errors.New("user's space brief has no value in 'userContactID' field")
	}

	return userSpaceBrief.UserContactID, nil
}
