package facade4invitus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/invitus/dbo4invitus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/random"
)

// CreateMassInviteRequest parameters for creating a mass invite
type CreateMassInviteRequest struct {
	Invite dbo4invitus.MassInviteDbo `json:"invite"`
}

// Validate validates parameters for creating a mass invite
func (request *CreateMassInviteRequest) Validate() error {
	return request.Invite.Validate()
}

// CreateMassInviteResponse creating a mass invite
type CreateMassInviteResponse struct {
	Invite MassInviteEntry
}

// CreateMassInvite creates a mass invite
func CreateMassInvite(ctx context.Context, userCtx facade.UserContext, request CreateMassInviteRequest) (response CreateMassInviteResponse, err error) {
	if err = request.Validate(); err != nil {
		err = fmt.Errorf("invalid request: %w", err)
		return
	}
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
	response CreateMassInviteResponse, err error,
) {
	invite := NewMassInviteEntryWithDbo(random.ID(7), &request.Invite)
	response.Invite = invite
	if err = tx.Insert(ctx, invite.Record); err != nil {
		return
	}
	return
}
