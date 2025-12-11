package facade4invitus

import (
	"context"
	"errors"
	"fmt"

	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/invitus/dbo4invitus"
	"github.com/sneat-co/sneat-go-core/facade"
)

// CreateMassInviteRequest parameters for creating a mass invite
type CreateMassInviteRequest struct {
	Invite dbo4invitus.InviteDbo `json:"invite"`
}

// Validate validates parameters for creating a mass invite
func (request *CreateMassInviteRequest) Validate() error {
	return request.Invite.Validate()
}

// CreateMassInvite creates a mass invite
func CreateMassInvite(ctx facade.ContextWithUser, request CreateMassInviteRequest) (response CreateInviteResponse, err error) {
	if err = request.Validate(); err != nil {
		err = fmt.Errorf("invalid request: %w", err)
		return
	}
	userCtx := ctx.User()
	if userCtx == nil || userCtx.GetUserID() == "" {
		err = errors.New("user context is required")
		return
	}
	var db dal.DB
	if db, err = facade.GetSneatDB(ctx); err != nil {
		return
	}
	err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		response, err = createMassInviteTx(ctx, tx, request)
		return
	})
	return
}

func createMassInviteTx(
	ctx context.Context, tx dal.ReadwriteTransaction, request CreateMassInviteRequest,
) (
	response CreateInviteResponse, err error,
) {
	invite := NewMassInviteEntryWithoutID(&request.Invite)
	if err = tx.Insert(ctx, invite.Record, dal.WithTimeStampStringID(dal.TimeStampAccuracyMicrosecond, 36, 3)); err != nil {
		return
	}
	invite.ID = invite.Record.Key().ID.(string)
	response.Invite = invite
	return
}
