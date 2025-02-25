package dal4contactus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

func GetContactByID(ctx context.Context, tx dal.ReadSession, spaceID coretypes.SpaceID, contactID string) (contact ContactEntry, err error) {
	contact = NewContactEntry(spaceID, contactID)
	return contact, GetContact(ctx, tx, contact)
}

func GetContact(ctx context.Context, tx dal.ReadSession, contact ContactEntry) error {
	return tx.Get(ctx, contact.Record)
}
