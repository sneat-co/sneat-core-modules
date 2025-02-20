package facade4invitus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/invitus/dbo4invitus"
)

// InvitesCollection table name
const InvitesCollection = "invites"

type PersonalInviteEntry = record.DataWithID[string, *dbo4invitus.PersonalInviteDbo]
type MassInviteEntry = record.DataWithID[string, *dbo4invitus.MassInviteDbo]

func NewPersonalInviteEntry(id string) (invite PersonalInviteEntry) {
	return NewPersonalInviteEntryWithDbo(id, new(dbo4invitus.PersonalInviteDbo))
}

func NewPersonalInviteEntryWithDbo(id string, dbo *dbo4invitus.PersonalInviteDbo) (invite PersonalInviteEntry) {
	invite.ID = id
	invite.Key = NewInviteKey(id)
	invite.Data = dbo
	invite.Record = dal.NewRecordWithData(invite.Key, invite.Data)
	return
}

func NewMassInviteEntry(id string) (invite MassInviteEntry) {
	return NewMassInviteEntryWithDbo(id, new(dbo4invitus.MassInviteDbo))
}

func NewMassInviteEntryWithDbo(id string, dbo *dbo4invitus.MassInviteDbo) (invite MassInviteEntry) {
	invite.ID = id
	invite.Key = NewInviteKey(id)
	invite.Data = dbo
	invite.Record = dal.NewRecordWithData(invite.Key, invite.Data)
	return
}
