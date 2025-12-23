package facade4invitus

import (
	"context"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/invitus/dbo4invitus"
)

// GetInviteByID returns an invitation record by ContactID
func GetInviteByID(ctx context.Context, getter dal.ReadSession, id string) (inviteDto *dbo4invitus.InviteDbo, inviteRecord dal.Record, err error) {
	inviteDto = new(dbo4invitus.InviteDbo)
	inviteRecord = dal.NewRecordWithData(NewInviteKey(id), inviteDto)
	return inviteDto, inviteRecord, getter.Get(ctx, inviteRecord)
}

// GetPersonalInviteByID returns an invitation record by ContactID
func GetPersonalInviteByID(ctx context.Context, getter dal.ReadSession, id string) (invite InviteEntry, err error) {
	invite = NewInviteEntry(id)
	return invite, getter.Get(ctx, invite.Record)
}

func GetInviteByInlineQueryID(ctx context.Context, getter dal.ReadSession, inlineQueryID string) (invite InviteEntry, err error) {
	q := dal.From(dal.NewCollectionRef(InvitesCollection, "", nil)).NewQuery().
		WhereField("inlineQueryID", dal.Equal, inlineQueryID).
		Limit(1).
		SelectIntoRecord(func() dal.Record {
			return NewInviteEntryWithDbo("", new(dbo4invitus.InviteDbo)).Record
		})
	var records []dal.Record
	if records, err = dal.ExecuteQueryAndReadAllToRecords(ctx, q, getter); err != nil {
		return
	}
	if len(records) == 0 {
		err = dal.ErrRecordNotFound
		return
	}
	r := records[0]
	key := r.Key()
	invite = InviteEntry{
		WithID: record.WithID[string]{
			ID:     key.ID.(string),
			Key:    key,
			Record: r,
		},
		Data: r.Data().(*dbo4invitus.InviteDbo),
	}
	return
}
