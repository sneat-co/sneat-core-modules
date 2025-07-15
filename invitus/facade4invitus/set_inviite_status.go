package facade4invitus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/invitus/dbo4invitus"
	"github.com/sneat-co/sneat-go-core/facade"
)

type inviteFields struct {
	inlineMessageID string
}

type InviteFieldArg func(v *inviteFields)

func WithInlineMessageID(inlineMessageID string) InviteFieldArg {
	return func(v *inviteFields) {
		v.inlineMessageID = inlineMessageID
	}
}

func SetInviteStatus(ctx context.Context, inviteID string, currentStatus, newStatus dbo4invitus.InviteStatus, fields ...InviteFieldArg) (invite InviteEntry, err error) {
	invite = NewInviteEntry(inviteID)
	var db dal.DB
	if db, err = facade.GetSneatDB(ctx); err != nil {
		return
	}
	args := new(inviteFields)
	for _, arg := range fields {
		arg(args)
	}

	if err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		if err = tx.Get(ctx, invite.Record); err != nil {
			return err
		}
		if currentStatus != "" && invite.Data.Status != currentStatus {
			return fmt.Errorf("invite status is %v, expected %v", invite.Data.Status, currentStatus)
		}
		if invite.Data.Status == newStatus {
			return nil
		}
		if invite.Data.Status != "" {
			switch newStatus {
			case dbo4invitus.InviteStatusSent:
				if invite.Data.Status != dbo4invitus.InviteStatusPending &&
					invite.Data.Status != dbo4invitus.InviteStatusSending {
					return fmt.Errorf(
						"only invite in status %s|%s can be moved to %s status, current invite status is: %s",
						dbo4invitus.InviteStatusPending, dbo4invitus.InviteStatusSending,
						dbo4invitus.InviteStatusSent,
						invite.Data.Status)
				}
			}
		}
		updates := []update.Update{update.ByFieldName("status", newStatus)}
		invite.Data.Status = newStatus
		if args.inlineMessageID != "" {
			invite.Data.InlineMessageID = args.inlineMessageID
			updates = append(updates, update.ByFieldName("inlineMessageID", args.inlineMessageID))
		}
		if err = invite.Data.Validate(); err != nil {
			return err
		}
		return tx.Update(ctx, invite.Key, updates)
	}); err != nil {
		return
	}
	return
}
