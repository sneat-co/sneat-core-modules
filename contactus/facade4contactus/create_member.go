package facade4contactus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/contactus/briefs4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-core-modules/contactus/models4contactus"
	"github.com/sneat-co/sneat-core-modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-core-modules/teamus/facade4teamus"
	"github.com/sneat-co/sneat-core-modules/userus/models4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"time"
)

// CreateMember adds members to a team
func CreateMember(
	ctx context.Context,
	user facade.User,
	request dal4contactus.CreateMemberRequest,
) (
	response dto4contactus.CreateContactResponse,
	err error,
) {
	createContactRequest := dto4contactus.CreateContactRequest{
		TeamRequest: request.TeamRequest,
		RelatedTo:   request.RelatedTo,
	}
	return CreateContact(ctx, user, createContactRequest)
}

func createMemberTx(
	ctx context.Context,
	tx dal.ReadwriteTransaction,
	user facade.User,
	request dal4contactus.CreateMemberRequest,
	params *dal4teamus.ModuleTeamWorkerParams[*models4contactus.ContactusTeamDto],
) (
	response dto4contactus.CreateContactResponse,
	err error,
) {
	now := time.Now()
	team := params.Team
	contactusTeam := dal4contactus.NewContactusTeamModuleEntry(params.Team.ID)
	if err := tx.Get(ctx, contactusTeam.Record); err != nil {
		return response, fmt.Errorf("failed to get contactus team record: %w", err)
	}

	if len(contactusTeam.Data.Contacts) == 0 {
		return response, errors.New("team has no members")
	}
	contactID, userMember := contactusTeam.Data.GetContactBriefByUserID(params.UserID)
	if userMember == nil {
		return response, errors.New("user does not belong to the team: " + params.UserID)
	}
	switch userMember.AgeGroup {
	case "", dbmodels.AgeGroupUnknown:
		if request.RelatedTo != nil {
			switch request.RelatedTo.RelatedAs {
			case dbmodels.RelationshipSpouse, dbmodels.RelationshipChild:
				userMember.AgeGroup = dbmodels.AgeGroupAdult
				userMemberKey := dal4contactus.NewContactKey(request.TeamID, contactID)
				if err = tx.Update(ctx, userMemberKey, []dal.Update{
					{
						Field: "ageGroup",
						Value: userMember.AgeGroup,
					},
				}); err != nil {
					return response, fmt.Errorf("failed to update member record: %w", err)
				}
			}
		}
	}
	memberBrief := request.ContactBrief

	if team.Data.Type == "family" {
		memberBrief.Roles = []string{
			const4contactus.TeamMemberRoleContributor,
		}
	}

	if memberBrief.Name.First != "" && briefs4contactus.IsUniqueShortTitle(memberBrief.Name.First, contactusTeam.Data.Contacts, const4contactus.TeamMemberRoleTeamMember) {
		memberBrief.ShortTitle = memberBrief.Name.First
	} else if memberBrief.Name.Nick != "" && briefs4contactus.IsUniqueShortTitle(memberBrief.Name.First, contactusTeam.Data.Contacts, const4contactus.TeamMemberRoleTeamMember) {
		memberBrief.ShortTitle = memberBrief.Name.Nick
	} else if memberBrief.Name.Full != "" {
		memberBrief.ShortTitle = getShortTitle(memberBrief.Name.Full, contactusTeam.Data.Contacts)
	} else if request.Title != "" {
		memberBrief.ShortTitle = getShortTitle(request.Title, contactusTeam.Data.Contacts)
	}

	//if request.Emails != "" {
	//	memberBrief.Avatar = &dbprofile.Avatar{
	//		Gravatar: fmt.Sprintf("%x", md5.Sum([]byte(strings.ToLower(request.Email)))),
	//	}
	//}

	//if memberBrief.Name.First != "" && memberBrief.Name.Last != "" {
	//
	//}
	contactID, err = dbmodels.GenerateIDFromNameOrRandom(memberBrief.Name, contactusTeam.Data.ContactIDs())
	if err != nil {
		return response, fmt.Errorf("failed to generate new member ItemID: %w", err)
	}

	var from string
	memberFoundByID := false
	for _, m := range contactusTeam.Data.Contacts {
		if m.UserID == params.UserID {
			memberFoundByID = true
			from = m.GetTitle()
			if from == "" {
				from = "userID=" + params.UserID
			}
		}
	}
	if !memberFoundByID {
		err = validation.NewErrBadRequestFieldValue("userID", "user does not belong to the team: userID="+params.UserID)
		return
	}
	if from == "" {
		err = validation.NewErrBadRequestFieldValue("userID", "team member has no title: userID="+params.UserID)
		return
	}
	{ // Update team record
		params.TeamUpdates = append(params.TeamUpdates,
			contactusTeam.Data.AddContact(contactID, &memberBrief),
		)
	}
	var contact dal4contactus.ContactEntry
	contact, err = facade4teamus.CreateMemberRecordFromBrief(ctx, tx, request.TeamID, contactID, memberBrief, now, params.UserID)
	if err != nil {
		return response, fmt.Errorf("failed to create member's record: %w", err)
	}

	if err = txUpdateMemberGroup(ctx, tx, params.Started, user.GetID(), params.Team.Data, params.Team.Key, params.TeamUpdates); err != nil {
		return response, fmt.Errorf("failed to update team record: %w", err)
	}
	params.TeamModuleEntry.Data.Contacts[contact.ID] = &contact.Data.ContactBrief
	if params.TeamModuleEntry.Record.Exists() {
		params.TeamModuleUpdates = append(params.TeamModuleUpdates, dal.Update{
			Field: "contacts." + contact.ID,
			Value: &contact.Data.ContactBrief,
		})
	} else {
		return response, errors.New("contactus team module entry does not exists")
	}
	if contactID == "" {
		panic("contactID is empty")
	}
	response.ID = contact.ID
	response.Data = contact.Data
	return
}

func getShortTitle(title string, members map[string]*briefs4contactus.ContactBrief) string {
	shortNames := models4userus.GetShortNames(title)
	for _, short := range shortNames {
		isUnique := true
		for _, m := range members {
			if m.ShortTitle == short.Name {
				isUnique = false
				break
			}
		}
		if isUnique {
			return short.Name
		}
	}
	return ""
}
