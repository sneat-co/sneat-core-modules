package facade4invitus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
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
