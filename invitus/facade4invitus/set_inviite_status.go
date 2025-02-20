package facade4invitus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/invitus/dbo4invitus"
	"github.com/sneat-co/sneat-go-core/facade"
)

func SetInviteStatus(ctx context.Context, inviteID string, currentStatus, newStatus dbo4invitus.InviteStatus) (invite InviteEntry, err error) {
	invite = NewInviteEntry(inviteID)
	var db dal.DB
	if db, err = facade.GetSneatDB(ctx); err != nil {
		return
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
			case dbo4invitus.InviteStatusActive:
				if invite.Data.Status != dbo4invitus.InviteStatusPending {
					return fmt.Errorf("only pending invites can be moved to active status, current invite status is: %s", invite.Data.Status)
				}
			}
		}
		invite.Data.Status = newStatus
		if err = invite.Data.Validate(); err != nil {
			return err
		}
		return tx.Update(ctx, invite.Key, []update.Update{update.ByFieldName("status", newStatus)})
	}); err != nil {
		return
	}
	return
}
