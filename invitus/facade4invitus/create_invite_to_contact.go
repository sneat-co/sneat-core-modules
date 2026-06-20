package facade4invitus

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/invitus/dbo4invitus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

// InviteContactRequest is a request DTO
type InviteContactRequest struct {
	ContactRequest
	//RemoteClient dbmodels.RemoteClientInfo `json:"remoteClient"`

	To    dbo4invitus.InviteTo `json:"to"`
	Roles []string             `json:"roles,omitempty"`
	//
	Send    bool   `json:"send,omitempty"`
	Message string `json:"message,omitempty"`
}

const maxMessageSize = 1000

// Validate returns error if not valid
func (v InviteContactRequest) Validate() error {
	if err := v.ContactRequest.Validate(); err != nil {
		return err
	}
	if v.To.ContactID != v.ContactID {
		return fmt.Errorf("contact ID in request does not match contact ID in 'to' field: %s != %s", v.ContactID, v.To.ContactID)
	}
	//if err := v.From.Validate(); err != nil {
	//	return validation.NewErrBadRequestFieldValue("from", err.Error())
	//}
	if err := v.To.Validate(); err != nil {
		return validation.NewErrBadRequestFieldValue("to", err.Error())
	}
	if len(v.Message) > maxMessageSize {
		return validation.NewErrBadRequestFieldValue("message", fmt.Sprintf("message length limit is %d characters max", maxMessageSize))
	}
	if v.To.Channel != "email" && v.Send {
		return fmt.Errorf("%w: at the moment invites can be sent only by email, channel='%s'", facade.ErrBadRequest, v.To.Channel)
	}
	return nil
}

// CreateOrReuseInviteToContact creates or reuses an invitation for a member
func CreateOrReuseInviteToContact(
	ctx facade.ContextWithUser,
	request InviteContactRequest,
	getRemoteClientInfo func() dbmodels.RemoteClientInfo,
) (
	response CreateInviteResponse,
	err error,
) {
	if ctx == nil {
		panic("argument 'ctx' cannot be nil")
	}
	if getRemoteClientInfo == nil {
		panic("argument 'getRemoteClientInfo' cannot be nil")
	}
	if err = request.Validate(); err != nil {
		err = validation.NewErrBadRequestFieldValue("request", fmt.Sprintf("invalid InviteContactRequest: %v", err))
		return
	}

	userID := ctx.User().GetUserID()

	err = contactusAccess.RunContactTx(ctx, request.ContactRequest,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, session ContactSession) (err error) {
			response.Space = session.Space()
			if err = session.GetRecords(ctx, tx); err != nil {
				return
			}
			if session.Space().Data.Type == coretypes.SpaceTypePrivate {
				return errors.New("private space does not support invites")
			}
			_, fromContactBrief := session.GetContactBriefByUserID(userID)
			if fromContactBrief == nil {
				return fmt.Errorf(
					"%w: current user does not belong to the space: user={id=%s}, space={id=%s,type=%s}",
					facade.ErrUnauthorized, userID, session.Space().ID, session.Space().Data.Type)
			}

			var invite InviteEntry

			//var inviteToContactBrief *dbo4contactus.InviteToContactBrief
			invite.ID = session.GetContactInviteBriefByChannelAndInviterUserID(request.To.Channel, userID)
			if invite.ID != "" {
				invite, err = GetPersonalInviteByID(ctx, tx, invite.ID)
				if invite.Data.Status == "active" || invite.Data.Status == "" {
					response.Invite = invite
					return
				}
				invite.Data = nil
				session.AppendContactDeleteInviteBrief(invite.ID)
				return
			}

			if invite.Data == nil {
				if invite, err = createPersonalInvite(ctx, tx, userID, request, session, getRemoteClientInfo); err != nil {
					return fmt.Errorf("failed to create personal invite: %w", err)
				}
			}
			response.Invite = invite
			return
		},
	)
	return
}

func createPersonalInvite(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	userID string,
	request InviteContactRequest,
	session ContactSession,
	getRemoteClientInfo func() dbmodels.RemoteClientInfo,
) (
	invite InviteEntry, err error,
) {
	if ctx == nil {
		panic("argument 'ctx' cannot be nil")
	}
	if tx == nil {
		panic("argument 'tx' cannot be nil")
	}
	if uid := strings.TrimSpace(userID); uid == "" {
		err = validation.NewErrRequestIsMissingRequiredField("userID")
		return
	} else if uid != userID {
		err = validation.NewErrBadRecordFieldValue("userID", "leading or trailing spaces")
		return
	}
	if session == nil {
		panic("argument 'session' cannot be nil")
	}
	if getRemoteClientInfo == nil {
		panic("argument 'getRemoteClientInfo' cannot be nil")
	}
	if request.SpaceID == "" {
		err = validation.NewErrRequestIsMissingRequiredField("spaceID")
		return
	}
	if request.To.ContactID == "" {
		err = validation.NewErrRequestIsMissingRequiredField("contactID")
		return
	}
	toContactID := session.Contacts()[request.To.ContactID]
	if toContactID == nil {
		err = errors.New("space has no 'to' contact with id=" + request.To.ContactID)
		return
	}
	request.To.Title = toContactID.GetTitle()

	fromContactID, fromBrief := session.GetContactBriefByUserID(userID)

	from := dbo4invitus.InviteFrom{
		InviteContact: dbo4invitus.InviteContact{
			UserID:    userID,
			ContactID: fromContactID,
			Title:     fromBrief.GetTitle(),
		},
	}
	to := request.To
	to.Title = toContactID.GetTitle()
	if !session.Space().Record.Exists() {
		err = fmt.Errorf("space record should not exist before creating a personal invite")
		return
	}
	inviteSpace := &dbo4invitus.InviteSpace{
		Type:  session.Space().Data.Type,
		Title: session.Space().Data.Title,
	}
	remoteClientInfo := getRemoteClientInfo()
	invite, err = createInviteToContactTx(
		ctx,
		tx,
		userID,
		remoteClientInfo,
		request.SpaceID,
		inviteSpace,
		from,
		to,
		!request.Send,
		request.Message,
		toContactID.Avatar,
	)
	if err != nil {
		err = fmt.Errorf("failed to create an invite record for a member: %w", err)
		return
	}
	if request.Send {
		if invite.Data.MessageID, err = sendInviteEmail(ctx, invite); err != nil {
			err = fmt.Errorf("%s: %w", FailedToSendEmail, err)
			return invite, err
		}
		if err = tx.Update(ctx, invite.Key,
			[]update.Update{update.ByFieldName("messageId", invite.Data.MessageID)}); err != nil {
			err = fmt.Errorf("failed to update invite record with message ContactID: %w", err)
			return
		}
	}

	session.AppendContactAddInviteBrief(invite.ID, userID, request.To.Channel, invite.Data.CreatedAt)
	return
}
