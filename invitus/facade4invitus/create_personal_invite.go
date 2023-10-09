package facade4invitus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-core-modules/invitus/models4invitus"
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

var randomInviteID = func() string {
	return random.ID(6)
}

var randomPinCode = func() string {
	return random.Digits(4)
}

// FailedToSendEmail error message
const FailedToSendEmail = "failed to send email"

func createInviteForMember(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	uid string,
	remoteClient dbmodels.RemoteClientInfo,
	team models4invitus.InviteTeam,
	from models4invitus.InviteFrom,
	to models4invitus.InviteToMember,
	composeOnly bool,
	inviterUserID,
	message string,
	toAvatar *dbprofile.Avatar,
) (id string, personalInvite *models4invitus.PersonalInviteDto, err error) {
	if err = team.Validate(); err != nil {
		err = fmt.Errorf("parameter 'team' is not valid: %w", err)
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
	teamID := team.ID
	if teamID == "" {
		err = validation.NewErrRecordIsMissingRequiredField("team.InviteID")
		return
	}
	team.ID = ""
	if team.Type == "family" && team.Title != "" {
		team.Title = ""
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
	personalInvite = &models4invitus.PersonalInviteDto{
		InviteDto: models4invitus.InviteDto{
			Status: "active",
			Pin:    randomPinCode(),
			TeamID: teamID,
			InviteBase: models4invitus.InviteBase{
				Type:    "personal",
				Channel: to.Channel,
				From:    from, // TODO: get user email
				To: &models4invitus.InviteTo{
					InviteContact: to.InviteContact,
				},
				ComposeOnly: composeOnly,
			},
			Created: dbmodels.CreatedInfo{
				At:     time.Now(),
				Client: remoteClient,
			},
			Team:    team,
			Message: message,
			Roles:   []string{"contributor"},
		},
		Address:        toAddressLower,
		ToTeamMemberID: briefs4contactus.GetFullContactID(teamID, to.MemberID),
		ToAvatar:       toAvatar,
	}
	id = randomInviteID()
	inviteKey := NewInviteKey(id)
	if err = personalInvite.Validate(); err != nil {
		err = fmt.Errorf("personal invite record data are not valid: %w", err)
		return
	}
	inviteRecord := dal.NewRecordWithData(inviteKey, personalInvite)
	if err = tx.Insert(ctx, inviteRecord); err != nil {
		err = fmt.Errorf("failed to insert a new invite record into database: %w", err)
		return
	}
	return
}
