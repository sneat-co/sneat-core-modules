package facade4invitus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/invitus/dbo4invitus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/slice"
	"github.com/strongo/validation"
	"slices"
	"time"
)

type ClaimInviteRequest struct {
	InviteRequest
	Operation InviteClaimOperation `json:"operation"`

	NoPinRequired bool `json:"noPinRequired,omitempty"`

	RemoteClient dbmodels.RemoteClientInfo `json:"remoteClient"`

	// TODO: Document why we need this and why it's called 'member'?
	//Member dbmodels.DtoWithID[*briefs4contactus.ContactBase] `json:"member"`

	//FullName string                      `json:"fullName"`
	//EmailAddress    string                      `json:"email"`
}

// Validate validates request
func (v *ClaimInviteRequest) Validate() error {
	if err := v.InviteRequest.Validate(); err != nil {
		return err
	}
	switch v.Operation {
	case "":
		return validation.NewErrRecordIsMissingRequiredField("operation")
	case InviteClaimOperationAccept, InviteClaimOperationDecline:
		// OK
	default:
		return validation.NewErrBadRequestFieldValue("operation", "invalid value: "+string(v.Operation))
	}
	//if err := v.Member.Validate(); err != nil {
	//	return validation.NewErrBadRequestFieldValue("member", err.Error())
	//}
	return nil
}

func ClaimInvite(ctx facade.ContextWithUser, r ClaimInviteRequest) (invite InviteEntry, err error) {
	invite = NewInviteEntry(r.InviteID)

	var db dal.DB

	if db, err = facade.GetSneatDB(ctx); err != nil {
		return
	}
	now := time.Now()
	userID := ctx.User().GetUserID()
	if err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		if err = tx.Get(ctx, invite.Record); err != nil {
			return
		}
		if invite.Data.Pin != r.Pin {
			return ErrInvitePinDoesNotMatch
		}
		if invite.Data.Type == dbo4invitus.InviteTypeMass {
			if invite.Data.Status == dbo4invitus.InviteStatusSending {
				invite.Data.Status = dbo4invitus.InviteStatusSent
				invite.Record.MarkAsChanged()
			}
			switch r.Operation {
			case InviteClaimOperationAccept:
				if !slices.Contains(invite.Data.AcceptedByUserIDs, userID) {
					invite.Data.AcceptedByUserIDs = append(invite.Data.AcceptedByUserIDs, userID)
					invite.Data.AcceptedCount += 1
					invite.Record.MarkAsChanged()
				}
				if slices.Contains(invite.Data.DeclinedByUserIDs, userID) {
					invite.Data.DeclinedByUserIDs = slice.RemoveInPlaceByValue(invite.Data.DeclinedByUserIDs, userID)
					invite.Data.DeclinedCount -= 1
					invite.Record.MarkAsChanged()
				}
			case InviteClaimOperationDecline:
				if !slices.Contains(invite.Data.DeclinedByUserIDs, userID) {
					invite.Data.DeclinedByUserIDs = append(invite.Data.DeclinedByUserIDs, userID)
					invite.Data.DeclinedCount += 1
					invite.Record.MarkAsChanged()
				}
				if slices.Contains(invite.Data.AcceptedByUserIDs, userID) {
					invite.Data.AcceptedByUserIDs = slice.RemoveInPlaceByValue(invite.Data.AcceptedByUserIDs, userID)
					invite.Data.DeclinedCount -= 1
					invite.Record.MarkAsChanged()
				}
			default:
				err = validation.NewErrBadRequestFieldValue("operation", "invalid value: "+string(r.Operation))
			}
		} else {
			switch r.Operation {
			case InviteClaimOperationAccept:
				invite.Data.Status = dbo4invitus.InviteStatusAccepted
				invite.Record.MarkAsChanged()
			case InviteClaimOperationDecline:
				invite.Data.Status = dbo4invitus.InviteStatusDeclined
				invite.Record.MarkAsChanged()
			default:
				err = validation.NewErrBadRequestFieldValue("operation", "invalid value: "+string(r.Operation))
			}
			invite.Data.Claimed = now
		}
		if invite.Record.HasChanged() {
			if err = tx.Set(ctx, invite.Record); err != nil {
				return
			}
		}
		return
	}); err != nil {
		return
	}
	return
}
