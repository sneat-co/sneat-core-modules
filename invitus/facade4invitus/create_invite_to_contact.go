package facade4invitus

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/update"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-core-modules/invitus/dbo4invitus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

// InviteContactRequest is a request DTO
type InviteContactRequest struct {
	dto4contactus.ContactRequest
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

	err = dal4contactus.RunContactWorker(ctx, request.ContactRequest,
		func(ctx facade.ContextWithUser, tx dal.ReadwriteTransaction, params *dal4contactus.ContactWorkerParams) (err error) {
			response.Contact = params.Contact
			response.ContactusSpace = params.SpaceModuleEntry
			response.Space = params.Space
			if err = tx.GetMulti(ctx, []dal.Record{
				params.Space.Record,
				params.Contact.Record,
				params.SpaceModuleEntry.Record,
			}); err != nil {
				return
			}
			if params.Space.Data.Type == coretypes.SpaceTypePrivate {
				return errors.New("private space does not support invites")
			}
			_, fromContactBrief := params.SpaceModuleEntry.Data.GetContactBriefByUserID(userID)
			if fromContactBrief == nil {
				return fmt.Errorf(
					"%w: current user does not belong to the space: user={id=%s}, space={id=%s,type=%s}",
					facade.ErrUnauthorized, userID, params.Space.ID, params.Space.Data.Type)
			}

			var invite InviteEntry

			//var inviteToContactBrief *dbo4contactus.InviteToContactBrief
			invite.ID, _ = params.Contact.Data.GetInviteBriefByChannelAndInviterUserID(request.To.Channel, userID)
			if invite.ID != "" {
				invite, err = GetPersonalInviteByID(ctx, tx, invite.ID)
				if invite.Data.Status == "active" || invite.Data.Status == "" {
					response.Invite = invite
					return
				}
				invite.Data = nil
				params.ContactUpdates = append(params.ContactUpdates, params.Contact.Data.DeleteInviteBrief(invite.ID))
				return
			}

			if invite.Data == nil {
				if invite, err = createPersonalInvite(ctx, tx, userID, request, params, getRemoteClientInfo); err != nil {
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
	params *dal4contactus.ContactWorkerParams,
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
	if params == nil {
		panic("argument 'params' cannot be nil")
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
	toContactID := params.SpaceModuleEntry.Data.Contacts[request.To.ContactID]
	if toContactID == nil {
		err = errors.New("space has no 'to' contact with id=" + request.To.ContactID)
		return
	}
	request.To.Title = toContactID.GetTitle()

	fromContactID, fromBrief := params.SpaceModuleEntry.Data.GetContactBriefByUserID(userID)

	from := dbo4invitus.InviteFrom{
		InviteContact: dbo4invitus.InviteContact{
			UserID:    userID,
			ContactID: fromContactID,
			Title:     fromBrief.GetTitle(),
		},
	}
	to := request.To
	to.Title = toContactID.GetTitle()
	if !params.Space.Record.Exists() {
		err = fmt.Errorf("space record should not exist before creating a personal invite")
		return
	}
	inviteSpace := &dbo4invitus.InviteSpace{
		Type:  params.Space.Data.Type,
		Title: params.Space.Data.Title,
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

	params.ContactUpdates = append(
		params.ContactUpdates,
		params.Contact.Data.AddInviteBrief(invite.ID, userID, request.To.Channel, invite.Data.CreatedAt),
	)
	return
}
