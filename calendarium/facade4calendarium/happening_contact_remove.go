package facade4calendarium

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/calendarium/dal4calendarium"
	"github.com/sneat-co/sneat-core-modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-core-modules/calendarium/models4calendarium"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

func RemoveParticipantFromHappening(ctx context.Context, user facade.User, request dto4calendarium.HappeningContactRequest) (err error) {
	if err = request.Validate(); err != nil {
		return
	}

	var worker = func(ctx context.Context, tx dal.ReadwriteTransaction, params *happeningWorkerParams) error {
		_, err := getHappeningContactRecords(ctx, tx, &request, params)
		if err != nil {
			return err
		}
		teamContactID := dbmodels.NewTeamItemID(request.Contact.TeamID, request.Contact.ID)
		switch params.Happening.Dto.Type {
		case "single":
			break // nothing to do
		case "recurring":
			var updates []dal.Update
			if updates, err = removeContactFromHappeningBriefInTeamDto(params.TeamModuleEntry, params.Happening, teamContactID); err != nil {
				return fmt.Errorf("failed to remove member from happening brief in team DTO: %w", err)
			}
			if len(updates) > 0 {
				params.TeamModuleUpdates = append(params.TeamModuleUpdates, updates...)
			}
		default:
			return fmt.Errorf("invalid happenning record: %w",
				validation.NewErrBadRecordFieldValue("type",
					fmt.Sprintf("unknown value: [%v]", params.Happening.Dto.Type)))
		}
		params.HappeningUpdates = append(params.HappeningUpdates, params.Happening.Dto.RemoveContact(request.Contact.TeamID, request.Contact.ID)...)
		params.HappeningUpdates = append(params.HappeningUpdates, params.Happening.Dto.RemoveParticipant(request.Contact.TeamID, request.Contact.ID)...)
		return err
	}

	if err = modifyHappening(ctx, user, request.HappeningRequest, worker); err != nil {
		return err
	}
	return nil
}

func removeContactFromHappeningBriefInTeamDto(
	calendariumTeam dal4calendarium.CalendariumTeamContext,
	happening models4calendarium.HappeningContext,
	teamContactID dbmodels.TeamItemID,
) (updates []dal.Update, err error) {
	happeningBrief := calendariumTeam.Data.GetRecurringHappeningBrief(happening.ID)
	if happeningBrief == nil {
		calendariumTeam.Data.RecurringHappenings[happening.ID] = &happening.Dto.HappeningBrief
	} else if _, ok := happeningBrief.Participants[string(teamContactID)]; ok {
		delete(happeningBrief.Participants, string(teamContactID))
		updates = append(updates, dal.Update{
			Field: fmt.Sprintf("recurringHappenings.%s.participants.%s", happening.ID, teamContactID),
			Value: dal.DeleteField,
		})
	}
	return updates, nil
}
