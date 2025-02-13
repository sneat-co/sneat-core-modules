package facade4invitus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-core-modules/invitus/dbo4invitus"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/sneat-co/sneat-go-core/models/dbprofile"
	"github.com/strongo/random"
	"github.com/strongo/validation"
	"net/mail"
	"strings"
	"time"
)

func NewInviteKey(inviteID string) *dal.Key {
	return dal.NewKeyWithID(InvitesCollection, inviteID)
}

var randomPinCode = func() string {
	return random.Digits(4)
}

// FailedToSendEmail error message
const FailedToSendEmail = "failed to send email"

func createInviteToContactTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	uid string,
	remoteClient dbmodels.RemoteClientInfo,
	space dbo4invitus.InviteSpace,
	from dbo4invitus.InviteFrom,
	to dbo4invitus.InviteTo,
	composeOnly bool,
	message string,
	toAvatar *dbprofile.Avatar,
) (invite PersonalInviteEntry, err error) {
	if err = space.Validate(); err != nil {
		err = fmt.Errorf("parameter 'space' is not valid: %w", err)
		return
	}
	if err = to.Validate(); err != nil {
		err = fmt.Errorf("parameter 'to' is not valid: %w", err)
		return
	}
	if err = from.Validate(); err != nil {
		err = fmt.Errorf("parameter 'from' is not valid: %w", err)
		return
	}
	spaceID := space.ID
	if spaceID == "" {
		err = validation.NewErrRecordIsMissingRequiredField("space.InviteID")
		return
	}
	space.ID = ""
	if space.Type == "family" && space.Title != "" {
		space.Title = ""
	}
	var toAddress *mail.Address
	if to.Address != "" {
		toAddress, err = mail.ParseAddress(to.Address)
		if err != nil {
			err = fmt.Errorf("failed to parse to.Address: %w", err)
			return
		}
	}
	var toAddressLower string
	if toAddress != nil {
		toAddressLower = strings.ToLower(toAddress.Address)
	}
	from.UserID = uid
	inviteDbo := &dbo4invitus.PersonalInviteDbo{
		InviteDbo: dbo4invitus.InviteDbo{
			Status:  "active",
			Pin:     randomPinCode(),
			SpaceID: spaceID,
			InviteBase: dbo4invitus.InviteBase{
				Type:    "personal",
				Channel: to.Channel,
				From:    from, // TODO: get user email
				To: &dbo4invitus.InviteTo{
					InviteContact: to.InviteContact,
				},
				ComposeOnly: composeOnly,
			},
			CreatedAt: time.Now(),
			Created: dbmodels.CreatedInfo{
				Client: remoteClient,
			},
			Space:   space,
			Message: message,
			Roles:   []string{"contributor"},
		},
		Address:          toAddressLower,
		ToSpaceContactID: briefs4contactus.GetFullContactID(spaceID, to.ContactID),
		ToAvatar:         toAvatar,
	}
	if err = inviteDbo.Validate(); err != nil {
		err = fmt.Errorf("personal invite record data are not valid: %w", err)
		return
	}

	var inviteKey *dal.Key
	if inviteKey, err = dal.NewKeyWithOptions(InvitesCollection, dal.WithRandomStringID(dal.RandomLength(6))); err != nil {
		return
	}
	inviteRecord := dal.NewRecordWithData(inviteKey, inviteDbo)
	//invite.ID = randomInviteID()
	//inviteKey := NewInviteKey(invite.ID)
	//inviteRecord := dal.NewRecordWithData(inviteKey, personalInvite)

	if err = tx.Insert(ctx, inviteRecord); err != nil {
		err = fmt.Errorf("failed to insert a new invite record into database: %w", err)
		return
	}
	invite = NewPersonalInviteEntryWithDto(inviteKey.ID.(string), inviteDbo)
	return
}
