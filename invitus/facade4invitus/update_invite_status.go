package facade4invitus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/invitus/dbo4invitus"
	"time"
)

func updateInviteStatus(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	uid string,
	now time.Time,
	invite InviteEntry,
	status dbo4invitus.InviteStatus,
) (err error) {
	var inviteUpdates []update.Update

	if invite.Data.Claimed != nil &&
		(status == dbo4invitus.InviteStatusSending || status == dbo4invitus.InviteStatusPending) {
		err = fmt.Errorf("claimed invite can not be moved to status %s", status)
		return err
	}

	if invite.Data.Expires != nil && invite.Data.Expires.Before(now) {
		err = fmt.Errorf("%w: expired at: %s", ErrInviteExpired, invite.Data.Expires)
		return
	}

	switch status {
	case dbo4invitus.InviteStatusPending:
		err = fmt.Errorf("not allowed to move invite to pending status")
		return err
	case dbo4invitus.InviteStatusSending:
		if invite.Data.Status != dbo4invitus.InviteStatusPending {
			err = fmt.Errorf("only pending invites can be moved to sending status, current invite status is: %s", invite.Data.Status)
			return err
		}
	case dbo4invitus.InviteStatusSent:
		switch invite.Data.Status {
		case dbo4invitus.InviteStatusPending, dbo4invitus.InviteStatusSending: // OK
		default:
			err = fmt.Errorf(
				"only invite in status %s|%s can be moved to %s status, current invite status is: %s",
				dbo4invitus.InviteStatusPending, dbo4invitus.InviteStatusSending,
				status,
				invite.Data.Status)
			return err
		}
	case dbo4invitus.InviteStatusAccepted:
		switch invite.Data.Status {
		case dbo4invitus.InviteStatusRevoked:
			err = fmt.Errorf("revoked invite can not be accepted: %w", ErrInviteIsRevoked)
			return
		case dbo4invitus.InviteStatusAccepted:
			return ErrInviteAlreadyAccepted
		}
	case dbo4invitus.InviteStatusDeclined:
		switch invite.Data.Status {
		case dbo4invitus.InviteStatusRevoked:
			err = fmt.Errorf("revoked invite can not be declined: %w", ErrInviteIsRevoked)
			return
		case dbo4invitus.InviteStatusDeclined:
			return // Nothing to do
		}
	case dbo4invitus.InviteStatusExpired:
		switch invite.Data.Status {
		case dbo4invitus.InviteStatusAccepted:
			err = fmt.Errorf("not allowed to expire an already claimed invite")
			return
		case dbo4invitus.InviteStatusDeclined:
			err = fmt.Errorf("not allowed to expire an already declined invite")
			return
		case dbo4invitus.InviteStatusRevoked:
			err = fmt.Errorf("not allowed to expire a revoked invite")
			return
		}
		return
	default:
		invite.Data.Claimed = &now
		inviteUpdates = append(inviteUpdates, update.ByFieldName("claimed", now))
	}
	invite.Data.Status = status
	invite.Data.To.UserID = uid
	inviteUpdates = append(inviteUpdates, update.ByFieldName("status", status))
	if invite.Data.To.UserID == "" && (invite.Data.Type == dbo4invitus.InviteTypePersonal || invite.Data.Type == dbo4invitus.InviteTypePrivate) {
		inviteUpdates = append(inviteUpdates, update.ByFieldName("to.userID", uid))
	}
	if err = invite.Data.Validate(); err != nil {
		return fmt.Errorf("personal invite record is not valid: %w", err)
	}
	if err = tx.Update(ctx, invite.Key, inviteUpdates); err != nil {
		return err
	}
	return err
}
