package facade4invitus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/invitus/dbo4invitus"
)

// InvitesCollection table name
const InvitesCollection = "invites"

type InviteEntry = record.DataWithID[string, *dbo4invitus.InviteDbo]

func NewInviteEntry(id string) (invite InviteEntry) {
	return NewInviteEntryWithDbo(id, new(dbo4invitus.InviteDbo))
}

func NewInviteEntryWithDbo(id string, dbo *dbo4invitus.InviteDbo) (invite InviteEntry) {
	invite.ID = id
	invite.Key = NewInviteKey(id)
	invite.Data = dbo
	invite.Record = dal.NewRecordWithData(invite.Key, invite.Data)
	return
}

func NewMassInviteEntry(id string) (invite InviteEntry) {
	return NewMassInviteEntryWithDbo(id, new(dbo4invitus.InviteDbo))
}

func NewMassInviteEntryWithDbo(id string, dbo *dbo4invitus.InviteDbo) (invite InviteEntry) {
	invite.ID = id
	invite.Key = NewInviteKey(id)
	invite.Data = dbo
	invite.Record = dal.NewRecordWithData(invite.Key, invite.Data)
	return
}
