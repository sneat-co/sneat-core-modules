package facade4invitus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/invitus/dbo4invitus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"strconv"
	"time"
)

// JoinInfoRequest request
type JoinInfoRequest struct {
	InviteID string `json:"inviteID"` // InviteDbo ContactID
	Pin      string `json:"pin"`
}

// Validate validates request
func (v *JoinInfoRequest) Validate() error {
	if v.InviteID == "" {
		return validation.NewErrRecordIsMissingRequiredField("id")
	}
	if v.Pin == "" {
		return validation.NewErrRequestIsMissingRequiredField("pin")
	}
	if _, err := strconv.Atoi(v.Pin); err != nil {
		return validation.NewErrBadRequestFieldValue("pin", "%pin is expected to be an integer")
	}
	return nil
}

type InviteInfo struct {
	Created time.Time                `json:"created"`
	Status  dbo4invitus.InviteStatus `json:"status"`
	From    dbo4invitus.InviteFrom   `json:"from"`
	To      *dbo4invitus.InviteTo    `json:"to"`
	Message string                   `json:"message,omitempty"`
}

func (v InviteInfo) Validate() error {
	if v.Status == "" {
		return validation.NewErrRecordIsMissingRequiredField("status")
	}
	if v.Created.IsZero() {
		return validation.NewErrRecordIsMissingRequiredField("created")
	}
	if err := v.From.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("from", err.Error())
	}
	if err := v.To.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("to", err.Error())
	}
	return nil
}

// JoinInfoResponse response
type JoinInfoResponse struct {
	SpaceID coretypes.SpaceID                                   `json:"spaceID"`
	Space   dbo4invitus.InviteSpace                             `json:"space"`
	Invite  InviteInfo                                          `json:"invite"`
	Member  *dbmodels.DtoWithID[*briefs4contactus.ContactBrief] `json:"member"`
}

func (v JoinInfoResponse) Validated() error {
	if err := v.Space.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("space", err.Error())
	}
	if err := v.Invite.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("space", err.Error())
	}
	if nil == v.Member {
		return validation.NewErrRecordIsMissingRequiredField("member")
	}
	if err := v.Member.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("member", err.Error())
	}
	return nil
}

// GetSpaceJoinInfo return join info
func GetSpaceJoinInfo(ctx context.Context, request JoinInfoRequest) (response JoinInfoResponse, err error) {
	if err = request.Validate(); err != nil {
		return
	}
	var db dal.DB
	if db, err = facade.GetSneatDB(ctx); err != nil {
		return
	}

	var inviteDbo *dbo4invitus.InviteDbo
	inviteDbo, _, err = GetInviteByID(ctx, db, request.InviteID)
	if err != nil {
		err = fmt.Errorf("failed to get invite record by ID=%s: %w", request.InviteID, err)
		return
	}
	if inviteDbo == nil {
		err = errors.New("invite record not found by ContactID: " + request.InviteID)
		return
	}

	if inviteDbo.Pin != request.Pin {
		err = fmt.Errorf("%v: %w",
			validation.NewErrBadRequestFieldValue("pin", "invalid pin"),
			facade.ErrForbidden,
		)
		return
	}
	var member dal4contactus.ContactEntry
	if inviteDbo.To.ContactID != "" {
		member = dal4contactus.NewContactEntry(inviteDbo.SpaceID, inviteDbo.To.ContactID)
		if err = db.Get(ctx, member.Record); err != nil {
			err = fmt.Errorf("failed to get space member's contact record: %w", err)
			return
		}
	}
	if inviteDbo.Space != nil {
		response.Space = *inviteDbo.Space

	}
	response.SpaceID = inviteDbo.SpaceID
	response.Invite.Status = inviteDbo.Status
	response.Invite.Created = inviteDbo.CreatedAt
	response.Invite.From = inviteDbo.From
	response.Invite.To = inviteDbo.To
	response.Invite.Message = inviteDbo.Message
	if inviteDbo.To.ContactID != "" {
		response.Member = &dbmodels.DtoWithID[*briefs4contactus.ContactBrief]{
			ID:   inviteDbo.To.ContactID,
			Data: &member.Data.ContactBrief,
		}
	}
	return response, nil
}
