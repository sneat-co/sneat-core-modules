package facade4schedulus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/schedulus/dal4schedulus"
	"github.com/sneat-co/sneat-core-modules/schedulus/dto4schedulus"
	"github.com/sneat-co/sneat-core-modules/schedulus/models4schedulus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
)

func RemoveParticipantFromHappening(ctx context.Context, user facade.User, request dto4schedulus.HappeningContactRequest) (err error) {
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
	schedulusTeam dal4schedulus.SchedulusTeamContext,
	happening models4schedulus.HappeningContext,
	teamContactID dbmodels.TeamItemID,
) (updates []dal.Update, err error) {
	happeningBrief := schedulusTeam.Data.GetRecurringHappeningBrief(happening.ID)
	if happeningBrief == nil {
		schedulusTeam.Data.RecurringHappenings[happening.ID] = &happening.Dto.HappeningBrief
	} else if _, ok := happeningBrief.Participants[string(teamContactID)]; ok {
		delete(happeningBrief.Participants, string(teamContactID))
		updates = append(updates, dal.Update{
			Field: fmt.Sprintf("recurringHappenings.%s.participants.%s", happening.ID, teamContactID),
			Value: dal.DeleteField,
		})
	}
	return updates, nil
}
