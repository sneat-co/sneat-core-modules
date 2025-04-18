package dal4contactus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
)

type ContactEntry = record.DataWithID[string, *dbo4contactus.ContactDbo]

func NewContactEntry(spaceID coretypes.SpaceID, contactID string) ContactEntry {
	return NewContactEntryWithData(spaceID, contactID, new(dbo4contactus.ContactDbo))
}

func NewContactEntryWithData(spaceID coretypes.SpaceID, contactID string, data *dbo4contactus.ContactDbo) (contact ContactEntry) {
	key := NewContactKey(spaceID, contactID)
	contact.ID = contactID
	contact.FullID = string(spaceID) + ":" + contactID
	contact.Key = key
	contact.Data = data
	contact.Record = dal.NewRecordWithData(key, data)
	return
}

func FindContactEntryByContactID(contacts []ContactEntry, contactID string) (contact ContactEntry, found bool) {
	for _, contact := range contacts {
		if contact.ID == contactID {
			return contact, true
		}
	}
	return contact, false
}

func GetContactusSpace(ctx context.Context, tx dal.ReadSession, contactusSpace ContactusSpaceEntry) (err error) {
	if tx == nil {
		if tx, err = facade.GetSneatDB(ctx); err != nil {
			return err
		}
	}
	return tx.Get(ctx, contactusSpace.Record)
}
