package cmds4anybot

import (
	"fmt"
	"github.com/bots-go-framework/bots-fw-store/botsfwmodels"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-core-modules/anybot"
	briefs4contactus2 "github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	dto4contactus2 "github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/facade4contactus"
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/strongoapp/person"
	"time"
)

var contactMessageCommand = botsfw.Command{
	Code: "contact_message",
	//Commands:   []string{"/ping"},
	InputTypes: []botinput.WebhookInputType{botinput.WebhookInputContact},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		now := time.Now()
		contactMessage := whc.Input().(botinput.WebhookContactMessage)

		ctx := whc.Context()

		chatData := whc.ChatData().(*anybot.SneatAppTgChatDbo)

		appUserID := whc.AppUserID()
		userCtx := facade.NewUserContext(appUserID)
		const userCanBeNonSpaceMember = false

		spaceID := chatData.GetSpaceRef().SpaceID()

		if spaceID == "" {
			var appUserData botsfwmodels.AppUserData
			if appUserData, err = whc.AppUserData(); err != nil {
				return
			}
			spaceID = appUserData.(*dbo4userus.UserDbo).GetFamilySpaceID()
			if spaceID == "" {
				m.Text = "Space is not set for the chat"
				return
			}
		}

		request := dto4contactus2.CreateContactRequest{
			Status: dbmodels.StatusActive,
			Type:   briefs4contactus2.ContactTypePerson, // TODO: Duplicate of request.Person.ContactBase.Status?
			SpaceRequest: dto4spaceus.SpaceRequest{
				SpaceID: spaceID,
			},
			Person: &dto4contactus2.CreatePersonRequest{
				ContactBase: briefs4contactus2.ContactBase{
					Status: dbmodels.StatusActive,
					ContactBrief: briefs4contactus2.ContactBrief{
						Type:     "person",
						Gender:   dbmodels.GenderUnknown,
						AgeGroup: dbmodels.AgeGroupUnknown,

						Names: &person.NameFields{
							FirstName: contactMessage.GetFirstName(),
							LastName:  contactMessage.GetLastName(),
						},
					},
				},
			},
		}
		request.Person.UpdatedBy = appUserID
		request.Person.UpdatedAt = now

		phoneNumber := contactMessage.GetPhoneNumber()
		if phoneNumber != "" {
			request.Person.Phones = append(request.Person.Phones, dbmodels.PersonPhone{
				Type:     "personal",
				Number:   phoneNumber,
				Verified: false,
				Note:     "From address book",
			})
		}

		//vCardStr := contactMessage.GetVCard()
		//if vCardStr != "" {
		//
		//}
		var response dto4contactus2.CreateContactResponse
		if response, err = facade4contactus.CreateContact(ctx, userCtx, userCanBeNonSpaceMember, request); err != nil {
			return
		}
		m.Text = fmt.Sprintf("New contact created: %s %s", response.Data.Names.FirstName, response.Data.Names.LastName)
		return
	},
}
