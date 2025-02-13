package facade4invitus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dbo4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-core-modules/invitus/dbo4invitus"
	"github.com/sneat-co/sneat-core-modules/spaceus/core4spaceus"
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
	ctx context.Context,
	userCtx facade.UserContext,
	request InviteContactRequest,
	getRemoteClientInfo func() dbmodels.RemoteClientInfo,
) (
	inviteBrief dbo4invitus.InviteBrief,
	contact dal4contactus.ContactEntry,
	err error,
) {
	if userCtx == nil || userCtx.GetUserID() == "" {
		err = errors.New("user context is required")
		return
	}
	if err = request.Validate(); err != nil {
		err = fmt.Errorf("invalid request: %w", err)
		return
	}
	err = dal4contactus.RunContactWorker(ctx, userCtx, request.ContactRequest,
		func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4contactus.ContactWorkerParams) (err error) {
			contact = params.Contact
			if err = tx.GetMulti(ctx, []dal.Record{
				params.Space.Record,
				params.Contact.Record,
				params.SpaceModuleEntry.Record,
			}); err != nil {
				return
			}
			if params.Space.Data.Type == core4spaceus.SpaceTypePrivate {
				return errors.New("private space does not support invites")
			}
			userID := params.UserID()
			fromContactID, fromContactBrief := params.SpaceModuleEntry.Data.GetContactBriefByUserID(userID)
			if fromContactBrief == nil {
				return fmt.Errorf(
					"%w: current user does not belong to the space: user={id=%s}, space={id=%s,type=%s}",
					facade.ErrUnauthorized, userID, params.Space.ID, params.Space.Data.Type)
			}

			var invite PersonalInviteEntry

			//var inviteToContactBrief *dbo4contactus.InviteToContactBrief
			invite.ID, _ = contact.Data.WithInvitesToContactBriefs.GetInviteBriefByChannelAndInviterUserID(request.To.Channel, userID)
			if invite.ID != "" {
				invite, err = GetPersonalInviteByID(ctx, tx, invite.ID)
				if invite.Data.Status == "active" {
					return
				}
				params.ContactUpdates = append(params.ContactUpdates, contact.Data.DeleteInviteBrief(invite.ID))
				return
			}

			if invite.Data == nil {
				fromMember := dal4contactus.NewContactEntry(request.SpaceID, fromContactID)

				if invite, err = createPersonalInvite(ctx, tx, userID, request, params, fromMember, getRemoteClientInfo); err != nil {
					return fmt.Errorf("failed to create personal invite record: %w", err)
				}
			}
			inviteBrief = dbo4invitus.NewInviteBriefFromDbo(invite.ID, invite.Data.InviteDbo)
			return
		},
	)
	return
}

func createPersonalInvite(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	uid string,
	request InviteContactRequest,
	params *dal4contactus.ContactWorkerParams,
	fromMember dal4contactus.ContactEntry,
	getRemoteClientInfo func() dbmodels.RemoteClientInfo,
) (
	invite PersonalInviteEntry, err error,
) {

	toContactID := params.SpaceModuleEntry.Data.Contacts[request.To.ContactID]
	if toContactID == nil {
		err = errors.New("space has no 'to' contact with id=" + request.To.ContactID)
		return
	}
	request.To.Title = toContactID.GetTitle()
	from := dbo4invitus.InviteFrom{
		InviteContact: dbo4invitus.InviteContact{
			UserID:    uid,
			ContactID: fromMember.ID,
			Title:     fromMember.Data.GetTitle(),
		},
	}
	to := request.To
	to.Title = toContactID.GetTitle()
	if !params.Space.Record.Exists() {
		err = fmt.Errorf("space record should not exist before creating a personal invite")
		return
	}
	inviteSpace := dbo4invitus.InviteSpace{
		ID:    request.SpaceID,
		Type:  params.Space.Data.Type,
		Title: params.Space.Data.Title,
	}
	remoteClientInfo := getRemoteClientInfo()
	invite, err = createInviteToContactTx(
		ctx,
		tx,
		uid,
		remoteClientInfo,
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
			[]dal.Update{
				{Field: "messageId", Value: invite.Data.MessageID},
			}); err != nil {
			err = fmt.Errorf("failed to update invite record with message ContactID: %w", err)
			return
		}
	}
	params.Contact.Data.Invites[invite.ID] = dbo4contactus.InviteToContactBrief{
		CreatedTime:     invite.Data.CreatedAt,
		CreatedByUserID: uid,
	}
	params.ContactUpdates = append(
		params.ContactUpdates,
		params.Contact.Data.AddInviteBrief(invite.ID, uid, request.To.Channel, invite.Data.CreatedAt))
	return
}
